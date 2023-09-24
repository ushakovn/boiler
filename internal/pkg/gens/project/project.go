package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ushakovn/boiler/internal/app/gen"
	desc "github.com/ushakovn/boiler/pkg/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

type project struct {
	workDir     string
	projectDir  string
	projectDesc *desc.Project
}

type Config struct {
	WorkDir    string
	ProjectDir string
}

func (c *Config) Validate() error {
	if c.WorkDir == "" {
		log.Printf("boilder: use default working directory path")
	}
	if c.ProjectDir == "" {
		log.Printf("boilder: use default project ")
		c.ProjectDir = "./config/project/config.textproto"
	}
	return nil
}

func NewProject(config Config) (gen.Generator, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config.Validate: %w", err)
	}
	return &project{
		workDir:    config.WorkDir,
		projectDir: config.ProjectDir,
	}, nil
}

func (g *project) Generate(ctx context.Context) error {
	if err := g.loadProjectDesc(); err != nil {
		return fmt.Errorf("g.loadProjectDesc: %w", err)
	}
	for _, dir := range g.projectDesc.Dirs {
		if err := g.genDirectory(ctx, dir, ""); err != nil {
			return fmt.Errorf("g.genDirectory %w", err)
		}
	}
	return nil
}

func (g *project) loadProjectDesc() error {
	buf, err := os.ReadFile(g.projectDir)
	if err != nil {
		return fmt.Errorf("os.ReadFile projectDir: %w", err)
	}
	proj := &desc.Project{}

	if err := prototext.Unmarshal(buf, proj); err != nil {
		return fmt.Errorf("prototext.Unmarshal: %w", err)
	}
	g.projectDesc = proj

	return nil
}

func (g *project) genFile(file *desc.File) error {
	path := fmt.Sprintf("%s.%s", file.Name, strings.TrimPrefix(file.Extension, "."))

	if _, err := os.Create(path); err != nil {
		return fmt.Errorf("os.CreateFile: %w", err)
	}
	if template := file.Template; template != nil {
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

func (g *project) genDirectory(ctx context.Context, dir *desc.Dir, parentPath string) error {
	path := g.buildFilepath(parentPath, dir.Name)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return fmt.Errorf("os.Mkdir dir: %w", err)
	}
	for _, file := range dir.Files {
		file.Name = g.buildFilepath(parentPath, dir.Name, file.Name)

		if err := g.genFile(file); err != nil {
			return fmt.Errorf("g.genFile file: %w", err)
		}
	}
	for _, nested := range dir.Nested {
		if err := g.genDirectory(ctx, nested, dir.Name); err != nil {
			return fmt.Errorf("g.genDirectory nested: %w", err)
		}
	}
	return nil
}

func (g *project) buildFilepath(parts ...string) string {
	pd := []string{g.workDir}
	pd = append(pd, parts...)
	p := filepath.Join(pd...)
	return p
}
