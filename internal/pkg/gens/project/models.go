package project

import (
	"os"
	"strings"
)

type projectDesc struct {
	Root *rootDesc `json:"root"`
}

type rootDesc struct {
	Files []*fileDesc      `json:"files"`
	Dirs  []*directoryDesc `json:"dirs"`
}

type directoryDesc struct {
	Dirs  []*directoryDesc `json:"dirs"`
	Name  *nameDesc        `json:"name"`
	Files []*fileDesc      `json:"files"`
}

type nameDesc struct {
	Value string `json:"value"`
	Env   string `json:"env"`
}

type fileDesc struct {
	// Path field it is combination of fileDesc.Name and path prefixes
	Path      string        `json:"-"`
	Name      string        `json:"name"`
	Extension string        `json:"extension"`
	Template  *templateDesc `json:"template"`
}

type templateDesc struct {
	Path       string `json:"path"`
	Executable bool   `json:"executable"`
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
