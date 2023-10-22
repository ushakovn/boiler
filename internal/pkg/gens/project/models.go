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
  Func  string `json:"func" yaml:"func"`
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

type execFunc func() any

func (g *project) setExecFunctions() {
  // Write other exec functions
  g.execFunctions = map[string]execFunc{
    "appName": func() any {
      return g.workDirFolder()
    },
  }
}

func (d *nameDesc) Execute(execFunctions map[string]execFunc) string {
  if d.Value != "" {
    return d.Value
  }
  if f, ok := execFunctions[d.Func]; ok {
    value, ok := f().(string)
    if ok {
      return value
    }
  }
  return ""
}
