package config

import (
  "fmt"

  "github.com/ushakovn/boiler/internal/pkg/stringer"
  "github.com/ushakovn/boiler/internal/pkg/validator"
)

type AppInfo struct {
  Name        string
  Version     string
  Description string
}

type configValues map[string]*configValue

type configParsed struct {
  App    configAppSection    `yaml:"app"`
  Custom configCustomSection `yaml:"custom"`
}

type configCustomSectionKey string

type configCustomSectionVal struct {
  Group       string `yaml:"group"`
  Type        string `yaml:"type"`
  Value       string `yaml:"value"`
  Description string `yaml:"description"`
}

type configCustomSection map[configCustomSectionKey]configCustomSectionVal

type configAppSection struct {
  Name        string `yaml:"name"`
  Version     string `yaml:"version"`
  Description string `yaml:"description"`
}

func (c configParsed) Validate() error {
  if err := c.App.Validate(); err != nil {
    return fmt.Errorf("invalid app section: %w", err)
  }
  if err := c.Custom.Validate(); err != nil {
    return fmt.Errorf("invalid custom section: %w", err)
  }
  return nil
}

func (c configAppSection) Validate() error {
  if c.Name == "" {
    return fmt.Errorf("name not specified")
  }
  if c.Version == "" {
    return fmt.Errorf("version not specified")
  }
  if c.Description == "" {
    return fmt.Errorf("description not specified")
  }
  return nil
}

func (c configCustomSection) Validate() error {
  fs := make([]validator.ValidateFunc, 0, 2*len(c))

  for key, value := range c {
    fs = append(fs, key.Validate)
    fs = append(fs, value.Validate(key.String()))
  }
  return validator.Validate(fs...)
}

func (c configCustomSectionKey) Validate() error {
  if stringer.IsWrongCase(c.String()) {
    return fmt.Errorf("invalid case for key: %s", c)
  }
  return nil
}

func (c configCustomSectionVal) Validate(sectionKey string) func() error {
  return func() error {
    var err error

    if c.Group == "" {
      err = fmt.Errorf("group not specified")
    }
    if c.Type == "" {
      err = fmt.Errorf("type not specified")
    }
    if c.Description == "" {
      err = fmt.Errorf("description not specified")
    }

    if err != nil {
      return fmt.Errorf("key %s: %w", sectionKey, err)
    }
    return nil
  }
}

func (c configCustomSectionKey) String() string {
  return string(c)
}
