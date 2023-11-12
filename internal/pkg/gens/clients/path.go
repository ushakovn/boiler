package clients

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

func parseGitHubProtoPath(protoPath string) *parsedGitHubProtoPath {
  // github.com/<owner>/<repo>/<path>.proto@<commit>
  pathParts := strings.SplitN(protoPath, "/", 4)

  if len(pathParts) != 4 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoPath)
  }
  ownerPart := pathParts[1]
  repoPart := pathParts[2]
  pathPart := pathParts[3]

  // <path>.proto@<commit>
  pathParts = strings.SplitN(pathPart, "@", 2)

  if len(pathParts) != 2 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoPath)
  }
  pathPart = pathParts[0]
  commitPart := pathParts[1]

  packageParts := strings.Split(pathPart, "/")

  if len(pathParts) == 0 {
    log.Fatalf("boiler: invalid path: %s. expected pattern: github.com/<owner>/<repo>/<path>.proto@<commit>", protoPath)
  }
  var packagePart string

  switch {
  // <path>.proto = <package>/<file>.proto
  case len(pathParts) >= 2:
    packagePart = packageParts[len(packageParts)-2]

  // <path>.proto = <file>.proto
  case len(pathParts) >= 1:
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
