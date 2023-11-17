package protodeps

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
  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
  "gopkg.in/yaml.v3"
)

type ProtoDeps struct {
  workDirPath  string
  goModuleName string

  githubClient  *resty.Client
  protoDepsPath string
}

type Config struct {
  GithubToken   string
  ProtoDepsPath string
}

func NewProtoDeps(config Config) (*ProtoDeps, error) {
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

  return &ProtoDeps{
    workDirPath:  workDirPath,
    goModuleName: goModuleName,

    githubClient:  githubClient,
    protoDepsPath: config.ProtoDepsPath,
  }, nil
}

func (g *ProtoDeps) Init(context.Context) error {
  if err := g.createVendorFolder(); err != nil {
    return fmt.Errorf("g.createVendorFolder: %w", err)
  }
  if err := g.createProtoDepsConfig(); err != nil {
    return fmt.Errorf("g.createProtoDepsConfig: %w", err)
  }
  return nil
}

func (g *ProtoDeps) createProtoDepsConfig() error {
  filePath := filepath.Join(g.workDirPath, "protodeps.yaml")

  if err := templater.ExecTemplateCopy(templates.ProtoDepsConfig, filePath, nil, nil); err != nil {
    return fmt.Errorf("execTemplateCopy: %w", err)
  }
  return nil
}

func (g *ProtoDeps) createVendorFolder() error {
  if _, err := filer.CreateNestedFolders(g.workDirPath, ".boiler", "vendor"); err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }
  return nil
}

func (g *ProtoDeps) Generate(ctx context.Context) error {
  protoDeps, err := g.collectProtoDeps()
  if err != nil {
    return fmt.Errorf("g.collectProtoDeps: %w", err)
  }

  if err = protoDeps.Validate(); err != nil {
    return fmt.Errorf("protoDeps.Validate: %w", err)
  }

  dstProtoFolder, err := filer.CreateNestedFolders(g.workDirPath, "pkg", "pb")
  if err != nil {
    return fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  if len(protoDeps.LocalDeps) > 0 {
    if err = g.generateLocalProtoDeps(ctx, dstProtoFolder, protoDeps.LocalDeps); err != nil {
      return fmt.Errorf("g.generateLocalProtoDeps: %w", err)
    }
  }

  if len(protoDeps.ExternalDeps) > 0 {
    if err = g.generateExternalProtoDeps(ctx, dstProtoFolder, protoDeps.ExternalDeps); err != nil {
      return fmt.Errorf("g.generateExternalProtoDeps: %w", err)
    }
  }

  return nil
}

func (g *ProtoDeps) generateExternalProtoDeps(ctx context.Context, dstProtoFolder string, protoDeps []*externalProtoDependency) error {
  localProtoDeps := make([]*localProtoDependency, 0, len(protoDeps))

  for _, externalProtoDep := range protoDeps {
    log.Infof("boiler: vendor external proto dependency: %s", externalProtoDep.Import)

    localProtoDep, err := g.vendorExternalProtoDep(ctx, externalProtoDep)
    if err != nil {
      return fmt.Errorf("g.vendorProtoDependency: %w", err)
    }
    localProtoDeps = append(localProtoDeps, localProtoDep)
  }
  if err := g.generateLocalProtoDeps(ctx, dstProtoFolder, localProtoDeps); err != nil {
    return fmt.Errorf("g.generateLocalProtoDeps: %w", err)
  }
  return nil
}

func (g *ProtoDeps) generateLocalProtoDeps(ctx context.Context, dstProtoFolder string, protoDeps []*localProtoDependency) error {
  for _, client := range protoDeps {
    if filer.IsExistedFile(client.Path) {
      continue
    }
    return fmt.Errorf("file does not exist: %s", client.Path)
  }

  protocOptions := []string{
    "--go_out=.",
    "--go-grpc_out=.",
  }
  var srcProtoPath []string

  for _, protoDep := range protoDeps {
    srcProtoPath = append(srcProtoPath, protoDep.Path)

    parsedPath := parseLocalProtoPath(protoDep.Path)

    dstProtoPath, err := filer.CreateNestedFolders(dstProtoFolder, parsedPath.Owner, parsedPath.Repo, parsedPath.Package)
    if err != nil {
      return fmt.Errorf("filer.CreateNestedFolders: %w", err)
    }
    dstRelProtoPath := strings.TrimPrefix(dstProtoPath, fmt.Sprint(g.workDirPath, "/"))

    options := []string{
      fmt.Sprint("--go_opt=M", protoDep.Path, "=", dstRelProtoPath),
      fmt.Sprint("--go-grpc_opt=M", protoDep.Path, "=", dstRelProtoPath),
    }

    protocOptions = append(protocOptions, options...)

    log.Infof("boiler: generate local proto dependency: %s", protoDep.Path)
  }
  protocOptions = append(protocOptions, srcProtoPath...)

  if err := executor.ExecCommandContext(ctx, "protoc", protocOptions...); err != nil {
    return fmt.Errorf("executor.ExecCommandContext: %w", err)
  }

  return nil
}

func (g *ProtoDeps) vendorExternalProtoDep(ctx context.Context, protoDep *externalProtoDependency) (*localProtoDependency, error) {
  loadedProto, err := g.loadExternalProtoDep(ctx, protoDep)
  if err != nil {
    return nil, fmt.Errorf("g.loadExternalProtoDep: %w", err)
  }
  parsedImport := parseGitHubProtoImport(protoDep.Import)

  nestedFolders := []string{".boiler", "vendor", parsedImport.Owner, parsedImport.Repo, parsedImport.Package}

  folderPath, err := filer.CreateNestedFolders(g.workDirPath, nestedFolders...)
  if err != nil {
    return nil, fmt.Errorf("filer.CreateNestedFolders: %w", err)
  }

  fileName := extractProtoFileName(protoDep.Import)
  filePath := filepath.Join(folderPath, fileName)

  if err = os.WriteFile(filePath, loadedProto, os.ModePerm); err != nil {
    return nil, fmt.Errorf("os.WriteFile: %w", err)
  }
  relFilePath := strings.TrimPrefix(filePath, fmt.Sprint(g.workDirPath, "/"))

  return &localProtoDependency{
    Path: relFilePath,
  }, nil
}

func (g *ProtoDeps) loadExternalProtoDep(ctx context.Context, protoDep *externalProtoDependency) ([]byte, error) {
  contentReq := buildGitHubContentRequest(protoDep.Import)

  contentResp := &githubContentResp{}
  var err error

  if err = g.doGithubRequestParsed(ctx, contentReq, contentResp); err != nil {
    return nil, fmt.Errorf("g.doGithubRequestParsed: %w", err)
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

func buildGitHubContentRequest(protoImport string) string {
  parsed := parseGitHubProtoImport(protoImport)

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

func (g *ProtoDeps) collectProtoDeps() (*protoDependencies, error) {
  fileBuf, err := os.ReadFile(g.protoDepsPath)
  if err != nil {
    return nil, fmt.Errorf("os.ReadFile: %w", err)
  }
  fileExt := filer.ExtractFileExtension(g.protoDepsPath)

  deps, err := parseProtoDeps(fileExt, fileBuf)
  if err != nil {
    return nil, fmt.Errorf("parseProtoDeps: %w", err)
  }
  return deps, nil
}

func parseProtoDeps(depsFileExt string, depsFileBuf []byte) (*protoDependencies, error) {
  var (
    deps *protoDependencies
    err  error
  )
  switch depsFileExt {
  case "yml", "yaml", "YML", "YAML":
    err = yaml.Unmarshal(depsFileBuf, &deps)
  case "json", "JSON":
    err = json.Unmarshal(depsFileBuf, &deps)
  default:
    err = fmt.Errorf("unsupported file extension: %s", depsFileExt)
  }
  return deps, err
}

func extractProtoFileName(protoImport string) string {
  const partsCount = 2
  fileName := filer.ExtractFileName(protoImport)

  if fileParts := strings.Split(fileName, "@"); len(fileParts) == partsCount {
    return fileParts[0]
  }
  return fileName
}
