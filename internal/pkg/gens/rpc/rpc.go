package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/ushakovn/boiler/internal/boiler/gen"
	"github.com/ushakovn/boiler/pkg/utils"
)

type rpc struct {
	workDirPath string
	rpcDirPath  string
	rpcDesc     *rootDesc
}

type Config struct {
	RpcDir string
}

func (c *Config) Validate() error {
	if c.RpcDir == "" {
		return fmt.Errorf("rpc directory not specfied")
	}
	return nil
}

func NewRpc(config Config) (gen.Generator, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config.Validate: %w", err)
	}
	pwd, err := utils.Env("PWD")
	if err != nil {
		return nil, err
	}
	workDirPath := pwd + "/______out______" // TODO: clear prefix

	return &rpc{
		rpcDirPath:  config.RpcDir,
		workDirPath: workDirPath,
	}, nil
}

func (g *rpc) Generate(context.Context) error {
	if err := g.loadRpcDesc(); err != nil {
		return fmt.Errorf("g.loadRpcDesc: %w", err)
	}
	if err := g.genRpcHandler(); err != nil {
		return fmt.Errorf("g.genRpcHandler: %w", err)
	}
	return nil
}

func (g *rpc) loadRpcDesc() error {
	buf, err := os.ReadFile(g.rpcDirPath)
	if err != nil {
		return fmt.Errorf("os.ReadFile projectDir: %w", err)
	}
	rpc := &rootDesc{}

	if err := json.Unmarshal(buf, rpc); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	g.rpcDesc = rpc

	return nil
}

func (g *rpc) genRpcHandler() error {
	const path = "./templates/rpc.template"

	t, err := template.ParseFiles(path)
	if err != nil {
		return fmt.Errorf("template.ParseFiles: %w", err)
	}
	f, err := os.Create("handler.go")
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	if err = t.Execute(f, g.rpcDesc); err != nil {
		return fmt.Errorf("t.Execute: %w", err)
	}
	return nil
}
