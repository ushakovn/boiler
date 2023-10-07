package project

import (
  "os"
  "strings"
)

type projectDesc struct {
  Root *rootDesc `json:"root"`
}

type rootDesc struct {
  Files []*fileDesc      `json:"files" yaml:"files"`
  Dirs  []*directoryDesc `json:"dirs" yaml:"dirs"`
}

type directoryDesc struct {
  Dirs  []*directoryDesc `json:"dirs" yaml:"dirs"`
  Name  *nameDesc        `json:"name" yaml:"name"`
  Files []*fileDesc      `json:"files" yaml:"files"`
}

type nameDesc struct {
  Value string `json:"value"`
  Env   string `json:"env"`
}

type fileDesc struct {
  // Path field it is combination of fileDesc.Name and path prefixes
  Path      string        `json:"-" yaml:"path"`
  Name      string        `json:"name" yaml:"name"`
  Extension string        `json:"extension" yaml:"extension"`
  Template  *templateDesc `json:"template" yaml:"template"`
}

type templateDesc struct {
  Path       string `json:"path" yaml:"path"`
  Compiled   string `json:"compiled" yaml:"compiled"`
  Executable bool   `json:"executable" yaml:"executable"`
}

func (d *nameDesc) String() string {
  if d.Value != "" {
    return d.Value
  }
  if d.Env != "" {
    env := os.Getenv(d.Env)

    if parts := strings.Split(env, `/`); len(parts) > 0 {
      return parts[len(parts)-1]
    }
    if parts := strings.Split(env, `\`); len(parts) > 0 {
      return parts[len(parts)-1]
    }
    return env
  }
  return ""
}
