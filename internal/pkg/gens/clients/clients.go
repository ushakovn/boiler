package clients

import (
  "context"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/go-resty/resty/v2"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/internal/pkg/executor"
  "github.com/ushakovn/boiler/internal/pkg/filer"
  "gopkg.in/yaml.v3"
)

type Clients struct {
  workDirPath  string
  goModuleName string

  githubClient    *resty.Client
  clientsDescPath string
}

type Config struct {
  GithubToken     string
  ClientsDescPath string
}

func (c *Config) Validate() error {
  if c.GithubToken == "" {
    return fmt.Errorf("github token not specified")
  }
  if c.ClientsDescPath == "" {
    return fmt.Errorf("clients path not specified")
  }
  return nil
}

func NewClients(config Config) (*Clients, error) {
  if err := config.Validate(); err != nil {
    return nil, err
  }
  workDirPath, err := filer.WorkDirPath()
  if err != nil {
    return nil, err
  }
  goModuleName, err := filer.ExtractGoModuleName(workDirPath)
  if err != nil {
    return nil, err
  }
  githubClient := resty.New()
  githubClient.SetAuthToken(config.GithubToken)

  return &Clients{
    workDirPath:  workDirPath,
    goModuleName: goModuleName,

    githubClient:    githubClient,
    clientsDescPath: config.ClientsDescPath,
  }, nil
}

func (g *Clients) Generate(ctx context.Context) error {
  clients, err := g.collectProtoClients()
  if err != nil {
    return fmt.Errorf("g.collectProtoClients: %w", err)
  }

  var (
    localClients  []*protoClient
    githubClients []*protoClient
  )

  for _, client := range clients.ProtoClients {
    log.Infof("boiler: validate proto client: %s", client.ProtoPath)

    if err = validateProtoClient(client); err != nil {
      return fmt.Errorf("validateProtoClient: %w", err)
    }
    if client.Local {
      localClients = append(localClients, client)
    } else {
      githubClients = append(githubClients, client)
    }
  }
  log.Infof("boiler: proto clients validated")

  dstProtoFolder, err := filer.CreateNestedFolders(g.workDirPath, "pkg", "pb")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  if len(localClients) > 0 {
    if err = g.generateLocalProtoClients(ctx, dstProtoFolder, localClients); err != nil {
      return fmt.Errorf("g.generateLocalProtoClients: %w", err)
    }
  }
  log.Infof("boiler: local proto clients generated")
  
  if len(githubClients) > 0 {
    if err = g.generateGithubProtoClients(ctx, dstProtoFolder, githubClients); err != nil {
      return fmt.Errorf("g.generateGithubProtoClients: %w", err)
    }
  }
  log.Infof("boiler: github proto clients generated")

  return nil
}

func (g *Clients) generateGithubProtoClients(ctx context.Context, dstProtoFolder string, clients []*protoClient) error {
  var err error

  for _, client := range clients {
    if client, err = g.vendorClientProto(ctx, client); err != nil {
      return fmt.Errorf("g.vendorClientProto: %w", err)
    }
  }
  if err = g.generateLocalProtoClients(ctx, dstProtoFolder, clients); err != nil {
    return fmt.Errorf("g.generateLocalProtoClients: %w", err)
  }
  return nil
}

func (g *Clients) generateLocalProtoClients(ctx context.Context, dstProtoFolder string, clients []*protoClient) error {
  for _, client := range clients {
    // Check file existence for specified proto path
    if filer.IsExistedFile(client.ProtoPath) {
      continue
    }
    return fmt.Errorf("file does not exist: %s", client.ProtoPath)
  }

  protocOptions := []string{
    "--go_out=.",
    "--go-grpc_out=.",
  }
  var srcProtoPath []string

  for _, client := range clients {
    srcProtoPath = append(srcProtoPath, client.ProtoPath)

    parsedPath := parseLocalProtoPath(client.ProtoPath)

    dstProtoPath, err := filer.CreateNestedFolders(dstProtoFolder, parsedPath.Owner, parsedPath.Repo, parsedPath.Package)
    if err != nil {
      return fmt.Errorf("filer.CreateNestedFolders: %w", err)
    }
    dstProtoPath = strings.TrimPrefix(dstProtoPath, fmt.Sprint(g.workDirPath, "/"))

    options := []string{
      fmt.Sprint("--go_opt=M", client.ProtoPath, "=", dstProtoPath),
      fmt.Sprint("--go-grpc_opt=M", client.ProtoPath, "=", dstProtoPath),
    }

    protocOptions = append(protocOptions, options...)

    log.Infof("boiler: generate proto client: %s", client.ProtoPath)
  }
  protocOptions = append(protocOptions, srcProtoPath...)

  if err := executor.ExecCommandContext(ctx, "protoc", protocOptions...); err != nil {
    return fmt.Errorf("executor.ExecCommandContext: %w", err)
  }

  return nil
}

func (g *Clients) vendorClientProto(ctx context.Context, client *protoClient) (*protoClient, error) {
  protoContent, err := g.pullClientProtoContent(ctx, client)
  if err != nil {
    return nil, fmt.Errorf("g.pullClientProtoContent: %w", err)
  }
  parsedPath := parseGitHubProtoPath(client.ProtoPath)

  nestedFolders := []string{".boiler", "vendor", parsedPath.Owner, parsedPath.Repo, parsedPath.Package}

  folderPath, err := filer.CreateNestedFolders(g.workDirPath, nestedFolders...)
  if err != nil {
    return nil, fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  fileName := extractProtoFileName(client.ProtoPath)
  filePath := filepath.Join(folderPath, fileName)

  if err = os.WriteFile(filePath, protoContent, os.ModePerm); err != nil {
    return nil, fmt.Errorf("os.WriteFile: %w", err)
  }
  // Set local file path for proto path
  client.ProtoPath = strings.TrimPrefix(filePath, fmt.Sprint(g.workDirPath, "/"))

  return client, nil
}

func (g *Clients) pullClientProtoContent(ctx context.Context, client *protoClient) ([]byte, error) {
  contentReq := buildGitHubContentRequest(client.ProtoPath)

  contentResp := &githubContentResp{}
  var err error

  if err = g.doGithubRequest(ctx, contentReq, contentResp); err != nil {
    return nil, fmt.Errorf("g.doGithubRequest: %w", err)
  }
  var content []byte

  switch {
  case contentResp.HasContent():
    if content, err = contentResp.DecodeContent(); err != nil {
      return nil, fmt.Errorf("contentResp.DecodeContent: %w", err)
    }
  case contentResp.HasDownloadUrl():
    if content, err = g.doGithubRequestRaw(ctx, contentResp.DownloadUrl); err != nil {
      return nil, fmt.Errorf("g.doGithubRequestRaw: %w", err)
    }
  }
  return content, nil
}

func buildGitHubContentRequest(protoPath string) string {
  parsed := parseGitHubProtoPath(protoPath)

  urlParts := []string{
    "https://api.github.com",
    "repos",
    parsed.Owner,
    parsed.Repo,
    "contents",
    parsed.Path,
  }
  contentUrl := strings.Join(urlParts, "/")

  contentUrl = fmt.Sprint(contentUrl, "?ref=", parsed.Commit)

  return contentUrl
}

func (g *Clients) collectProtoClients() (*protoClients, error) {
  fileBuf, err := os.ReadFile(g.clientsDescPath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile: %w", err)
  }
  fileExt := filer.ExtractFileExtension(g.clientsDescPath)

  clients, err := parseProtoClients(fileExt, fileBuf)
  if err != nil {
    return nil, fmt.Errorf("parseProtoClients: %w", err)
  }
  return clients, nil
}

func parseProtoClients(fileExtension string, buf []byte) (*protoClients, error) {
  var (
    desc *protoClients
    err  error
  )
  switch fileExtension {
  case "yml", "yaml", "YML", "YAML":
    err = yaml.Unmarshal(buf, &desc)
  case "json", "JSON":
    err = json.Unmarshal(buf, &desc)
  default:
    err = fmt.Errorf("unsupported file extension: %s", fileExtension)
  }
  return desc, err
}

func extractProtoFileName(protoPath string) string {
  const partsCount = 2
  fileName := filer.ExtractFileName(protoPath)

  if fileParts := strings.Split(fileName, "@"); len(fileParts) == partsCount {
    return fileParts[0]
  }
  return fileName
}
