package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ushakovn/boiler/internal/boiler/gen"
	"github.com/ushakovn/boiler/pkg/utils"
)

type project struct {
	projectCfgPath string
	workDirPath    string
	projectDesc    *projectDesc
}

type Config struct{}

func NewProject(_ Config) (gen.Generator, error) {
	pwd, err := utils.Env("PWD")
	if err != nil {
		return nil, err
	}
	workDirPath := pwd + "/______out______" // TODO: clear prefix

	gopath, err := utils.Env("GOPATH")
	if err != nil {
		return nil, err
	}
	projectConfigPath := filepath.Join(gopath, "/config/boiler/project.json")

	return &project{
		projectCfgPath: projectConfigPath,
		workDirPath:    workDirPath,
	}, nil
}

func (g *project) Generate(ctx context.Context) error {
	if err := g.loadProjectDesc(); err != nil {
		return fmt.Errorf("g.loadProjectDesc: %w", err)
	}
	for _, file := range g.projectDesc.Root.Files {
		file.Path = g.buildPath(file.Name)

		if err := g.genFile(file); err != nil {
			return fmt.Errorf("g.genFile: %w", err)
		}
	}
	for _, dir := range g.projectDesc.Root.Dirs {
		if err := g.genDirectory(ctx, dir, ""); err != nil {
			return fmt.Errorf("g.genDirectory %w", err)
		}
	}
	return nil
}

func (g *project) loadProjectDesc() error {
	buf, err := os.ReadFile(g.projectCfgPath)
	if err != nil {
		return fmt.Errorf("os.ReadFile projectDir: %w", err)
	}
	proj := &projectDesc{}

	if err = json.Unmarshal(buf, proj); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	g.projectDesc = proj

	return nil
}

func (g *project) genFile(file *fileDesc) error {
	var path string

	if path = file.Path; path == "" {
		return fmt.Errorf("file.Path not specified")
	}
	if extension := file.Extension; extension != "" {
		extension = strings.TrimPrefix(file.Extension, ".")
		path = fmt.Sprintf("%s.%s", file.Path, extension)
	}
	if _, err := os.Create(path); err != nil {
		return fmt.Errorf("os.CreateFile: %w", err)
	}
	if template := file.Template; template != nil {
		// Copy content from template to new created file
		if !template.Executable {
			buf, err := os.ReadFile(template.Path)
			if err != nil {
				return fmt.Errorf("os.ReadFile: %w", err)
			}
			if err := os.WriteFile(path, buf, os.ModePerm); err != nil {
				return fmt.Errorf("os.WriteFile: %w", err)
			}
		}
	}

	return nil
}

func (g *project) genDirectory(ctx context.Context, dir *directoryDesc, parentPath string) error {
	path := g.buildPath(parentPath, dir.Name.String())

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return fmt.Errorf("os.Mkdir dir: %w", err)
	}
	for _, file := range dir.Files {
		file.Path = g.buildPath(parentPath, dir.Name.String(), file.Name)

		if err := g.genFile(file); err != nil {
			return fmt.Errorf("g.genFile file: %w", err)
		}
	}
	for _, nested := range dir.Dirs {
		if err := g.genDirectory(ctx, nested, dir.Name.String()); err != nil {
			return fmt.Errorf("g.genDirectory nested: %w", err)
		}
	}
	return nil
}

func (g *project) buildPath(parts ...string) string {
	pd := []string{g.workDirPath}
	pd = append(pd, parts...)
	p := filepath.Join(pd...)
	return p
}
