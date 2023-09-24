package rpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/ushakovn/boiler/internal/app/gen"
	desc "github.com/ushakovn/boiler/pkg/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

type rpc struct {
	workDir string
	rpcDir  string
	rpcDesc *desc.Rpc
}

type Config struct {
	WorkDir string
	RpcDir  string
}

func (c *Config) Validate() error {
	if c.WorkDir == "" {
		log.Printf("boilder: use default working directory path")
	}
	if c.RpcDir == "" {
		return fmt.Errorf("rpc directory not specfied")
	}
	return nil
}

func NewRpc(config Config) (gen.Generator, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config.Validate: %w", err)
	}
	return &rpc{
		workDir: config.WorkDir,
		rpcDir:  config.RpcDir,
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
	buf, err := os.ReadFile(g.rpcDir)
	if err != nil {
		return fmt.Errorf("os.ReadFile projectDir: %w", err)
	}
	proj := &desc.Rpc{}

	if err := prototext.Unmarshal(buf, proj); err != nil {
		return fmt.Errorf("prototext.Unmarshal: %w", err)
	}
	g.rpcDesc = proj

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
