package protodeps

import (
  "errors"
  "fmt"
  "strings"

  "github.com/samber/lo"
  "github.com/ushakovn/boiler/internal/pkg/builder"
  "github.com/ushakovn/boiler/internal/pkg/validator"
)

type protoDependencies struct {
  AppDeps      []*externalProtoDependency `yaml:"app_deps" json:"app_deps"`
  LocalDeps    []*localProtoDependency    `yaml:"local_deps" json:"local_deps"`
  ExternalDeps []*externalProtoDependency `yaml:"external_deps" json:"external_deps"`
}

type localProtoDependency struct {
  Path string `yaml:"path" json:"path"`
}

type externalProtoDependency struct {
  Import string `yaml:"import" json:"import"`
}

func (d *protoDependencies) Validate() error {
  fs := make([]validator.ValidateFunc, 0, len(d.LocalDeps)+len(d.ExternalDeps))

  for _, protoDep := range d.AppDeps {
    fs = append(fs, protoDep.Validate)
  }
  for _, protoDep := range d.LocalDeps {
    fs = append(fs, protoDep.Validate)
  }
  for _, protoDep := range d.ExternalDeps {
    fs = append(fs, protoDep.Validate)
  }
  return validator.Validate(fs...)
}

func (d *localProtoDependency) Validate() error {
  const pathPrefix = "proto"

  if !strings.HasPrefix(d.Path, pathPrefix) {
    return fmt.Errorf("invalid local proto dependency %s: must be placed in %s", d.Path, pathPrefix)
  }
  return nil
}

func (d *externalProtoDependency) Validate() error {
  partsImport := strings.Split(d.Import, "@")

  if len(partsImport) != 2 {
    return fmt.Errorf("invalid github proto dependency: %s", d.Import)
  }
  protoPath := partsImport[0]

  const pathPrefix = "github.com"

  if !strings.HasPrefix(protoPath, pathPrefix) {
    return fmt.Errorf("invalid github proto dependency %s: must be placed in %s", d.Import, pathPrefix)
  }
  return nil
}

func newProtoDependencies() *protoDependencies {
  return &protoDependencies{}
}

func (d *protoDependencies) HasAppDeps() bool {
  return len(d.AppDeps) != 0
}

func (d *protoDependencies) HasLocalDeps() bool {
  return len(d.LocalDeps) != 0
}

func (d *protoDependencies) HasExternalDeps() bool {
  return len(d.ExternalDeps) != 0
}

func (d *protoDependencies) HasDeps() bool {
  return d.HasAppDeps() || d.HasLocalDeps() || d.HasExternalDeps()
}

func (d *protoDependencies) checkDuplicates() error {
  type namedProtoDeps struct {
    name string
    deps any
  }
  protoDeps := []*namedProtoDeps{
    {name: "app", deps: d.AppDeps},
    {name: "external", deps: d.ExternalDeps},
    {name: "local", deps: d.LocalDeps},
  }
  b := builder.NewBuilder()

  for _, protoDep := range protoDeps {
    if dup := findProtoDepsDuplicates(protoDep.deps); len(dup) > 0 {
      b.Write("\t%s proto deps duplicates: %v\n", protoDep.name, dup)
    }
  }
  if b.Count() != 0 {
    return errors.New(b.String())
  }
  return nil
}

func findProtoDepsDuplicates(protoDeps any) []string {
  paths := collectProtoPaths(protoDeps)
  return lo.FindDuplicates(paths)
}

func collectProtoPaths(protoDeps any) []string {
  var paths []string

  switch deps := protoDeps.(type) {
  case []*externalProtoDependency:
    paths = lo.Map(deps, func(dep *externalProtoDependency, _ int) string {
      return dep.Import
    })
  case []*localProtoDependency:
    paths = lo.Map(deps, func(dep *localProtoDependency, _ int) string {
      return dep.Path
    })
  }
  return paths
}

func filterExternalProtoDeps(protoDeps, protoDepsDump []*externalProtoDependency) []*externalProtoDependency {
  paths := filterProtoDepsPaths(
    collectProtoPaths(protoDeps),
    collectProtoPaths(protoDepsDump),
  )
  return lo.Map(paths, func(path string, _ int) *externalProtoDependency {
    return &externalProtoDependency{Import: path}
  })
}

func filterLocalProtoDeps(protoDeps, protoDepsDump []*localProtoDependency) []*localProtoDependency {
  paths := filterProtoDepsPaths(
    collectProtoPaths(protoDeps),
    collectProtoPaths(protoDepsDump),
  )
  return lo.Map(paths, func(path string, _ int) *localProtoDependency {
    return &localProtoDependency{Path: path}
  })
}

func filterProtoDepsPaths(protoDeps, protoDepsDump []string) []string {
  m := map[string]struct{}{}

  for _, path := range protoDepsDump {
    m[path] = struct{}{}
  }
  var paths []string

  for _, path := range protoDeps {
    if _, ok := m[path]; ok {
      continue
    }
    m[path] = struct{}{}
    paths = append(paths, path)
  }
  return paths
}
