package protodeps

import (
  "fmt"
  "path/filepath"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/validator"
)

type protoDependencies struct {
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
  validators := make([]validator.Validator, 0, len(d.LocalDeps)+len(d.ExternalDeps))

  for _, protoDep := range d.LocalDeps {
    validators = append(validators, protoDep)
  }
  for _, protoDep := range d.ExternalDeps {
    validators = append(validators, protoDep)
  }
  return validator.Validate(validators...)
}

func (d *localProtoDependency) Validate() error {
  pathPrefix := filepath.Join(".boiler", "vendor")

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
  pathPrefix := "github.com"

  if !strings.HasPrefix(protoPath, pathPrefix) {
    return fmt.Errorf("invalid github proto dependency %s: must be placed in %s", d.Import, pathPrefix)
  }
  return nil
}

func newProtoDependencies() *protoDependencies {
  return &protoDependencies{}
}

func (d *protoDependencies) HasLocalDeps() bool {
  return len(d.LocalDeps) != 0
}

func (d *protoDependencies) HasExternalDeps() bool {
  return len(d.ExternalDeps) != 0
}

func (d *protoDependencies) HasDeps() bool {
  return d.HasLocalDeps() || d.HasExternalDeps()
}
