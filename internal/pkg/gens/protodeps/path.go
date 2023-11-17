package protodeps

import (
  "fmt"
  "path/filepath"
  "strings"

  log "github.com/sirupsen/logrus"
)

type parsedLocalProtoPath struct {
  Owner   string
  Repo    string
  Package string
}

func parseLocalProtoPath(protoPath string) *parsedLocalProtoPath {
  pathPrefix := filepath.Join(".boiler", "vendor")

  protoPath = strings.TrimPrefix(protoPath, fmt.Sprint(pathPrefix, "/"))
  pathParts := strings.Split(protoPath, "/")

  var (
    ownerPart   string
    repoPart    string
    packagePart string
  )

  switch {
  // <owner>/<repo>/<path>.proto
  // <path>.proto = .../.../<package>/<file>.proto
  case len(pathParts) >= 4 && strings.HasSuffix(pathParts[len(pathParts)-1], ".proto"):
    ownerPart = pathParts[0]
    repoPart = pathParts[1]
    packagePart = pathParts[len(pathParts)-2]

  // <repo>/<path>.proto
  // <path>.proto = .../.../<package>/<file>.proto
  case len(pathParts) >= 3 && strings.HasSuffix(pathParts[len(pathParts)-1], ".proto"):
    repoPart = pathParts[1]
    packagePart = pathParts[len(pathParts)-2]

  // <path>.proto
  // <path>.proto = .../.../<package>/<file>.proto
  case len(pathParts) >= 2 && strings.HasSuffix(pathParts[len(pathParts)-1], ".proto"):
    packagePart = pathParts[len(pathParts)-2]
  }

  return &parsedLocalProtoPath{
    Owner:   ownerPart,
    Repo:    repoPart,
    Package: packagePart,
  }
}

type parsedGitHubProtoPath struct {
  Owner   string
  Repo    string
  Path    string
  Package string
  Commit  string
}

func parseGitHubProtoImport(protoImport string) *parsedGitHubProtoPath {
  // github.com/<owner>/<repo>/<path>.proto@<commit>
  partsImport := strings.SplitN(protoImport, "/", 4)

  if len(partsImport) != 4 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoImport)
  }
  ownerPart := partsImport[1]
  repoPart := partsImport[2]
  pathPart := partsImport[3]

  // <path>.proto@<commit>
  partsImport = strings.SplitN(pathPart, "@", 2)

  if len(partsImport) != 2 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoImport)
  }
  pathPart = partsImport[0]
  commitPart := partsImport[1]

  packageParts := strings.Split(pathPart, "/")

  if len(partsImport) == 0 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoImport)
  }
  var packagePart string

  switch {
  // <path>.proto = <package>/<file>.proto
  case len(partsImport) >= 2:
    packagePart = packageParts[len(packageParts)-2]

  // <path>.proto = <file>.proto
  case len(partsImport) >= 1:
    packagePart = packageParts[len(packageParts)-1]
    packagePart = strings.TrimSuffix(packagePart, ".proto")
  }

  return &parsedGitHubProtoPath{
    Owner:   ownerPart,
    Repo:    repoPart,
    Path:    pathPart,
    Package: packagePart,
    Commit:  commitPart,
  }
}
