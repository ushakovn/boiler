package config

import (
  "errors"
  "fmt"
  "strings"

  "github.com/ushakovn/boiler/internal/pkg/builder"
  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/validator"
  "github.com/ushakovn/boiler/pkg/config/types"
)

type AppInfo struct {
  Name        string
  Version     string
  Description string
}

type configValues map[string]types.Value

type Parsed struct {
  App    *AppSection   `yaml:"app"`
  Custom CustomSection `yaml:"custom"`
}

type CustomSectionKey string

type CustomSectionVal struct {
  Group       string `yaml:"group"`
  Type        string `yaml:"type"`
  Value       string `yaml:"value"`
  Description string `yaml:"description"`
}

type CustomSection map[CustomSectionKey]*CustomSectionVal

type AppSection struct {
  Name        string `yaml:"name"`
  Version     string `yaml:"version"`
  Description string `yaml:"description"`
}

func (c *Parsed) Validate() error {
  b := builder.NewBuilder()

  if err := c.App.Validate(); err != nil {
    b.Write("invalid app section:\n%v\n", err)
  }
  if err := c.Custom.Validate(); err != nil {
    b.Write("invalid custom section:\n\n%v", err)
  }
  if s := b.String(); s != "" {
    return errors.New(s)
  }
  return nil
}

func (c *AppSection) Validate() error {
  b := builder.NewBuilder()

  if c.Name == "" {
    b.Write("\tname not specified\n")
  }
  if c.Version == "" {
    b.Write("\tversion not specified\n")
  }
  if c.Description == "" {
    b.Write("\tdescription not specified\n")
  }
  if s := b.String(); s != "" {
    return errors.New(s)
  }
  return nil
}

func (c CustomSection) Validate() error {
  keys := map[CustomSectionKey]struct{}{}

  for key := range c {
    if _, ok := keys[key]; !ok {
      keys[key] = struct{}{}
      continue
    }
    return fmt.Errorf("custom section contains duplicated key: %s", key)
  }
  fs := make([]validator.ValidateFunc, 0, 2*len(c))

  for key, value := range c {
    fs = append(fs, value.Validate(key.String()))
    fs = append(fs, key.Validate)
  }
  return validator.Validate(fs...)
}

func (c CustomSectionKey) Validate() error {
  if stringer.IsWrongCase(c.String()) {
    return fmt.Errorf("invalid case for key: %s\n", c)
  }
  return nil
}

func (c *CustomSectionVal) Validate(sectionKey string) func() error {
  return func() error {
    b := builder.NewBuilder()

    if c.Group == "" {
      b.Write("\tgroup not specified\n")
    }
    if c.Type == "" {
      b.Write("\ttype not specified\n")
    }
    if !types.IsValid(c.Type) && c.Type != "" {
      b.Write("\tinvalid value type: %s\n", c.Type)
    }
    if c.Description == "" {
      b.Write("\tdescription not specified\n")
    }
    if !strings.HasPrefix(sectionKey, c.Group) {
      b.Write("\tnot contains group prefix: %s\n", c.Group)
    }
    if b.Count() != 0 {
      return fmt.Errorf("custom section: %s invalid:\n%s", sectionKey, b.String())
    }
    return nil
  }
}

func (c CustomSectionKey) String() string {
  return string(c)
}
