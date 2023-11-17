package project

type projectDesc struct {
  Root *rootDesc `json:"root" yaml:"root"`
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
  Value string `json:"value" yaml:"value"`
}

type fileDesc struct {
  // Path field it is combination of fileDesc.Name and path prefixes
  Path      string        `json:"-" yaml:"-"`
  Name      string        `json:"name" yaml:"name"`
  Extension string        `json:"extension" yaml:"extension"`
  Template  *templateDesc `json:"template" yaml:"template"`
}

type templateDesc struct {
  Name string `json:"name" yaml:"name"`
}
