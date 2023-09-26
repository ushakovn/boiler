package factory

import (
	"context"
	"fmt"

	"github.com/ushakovn/boiler/internal/boiler/gen"
	"github.com/ushakovn/boiler/internal/pkg/gens/project"
	"github.com/ushakovn/boiler/internal/pkg/gens/rpc"
)

type Typ string

const (
	ProjectTyp Typ = "project"
	RpcType    Typ = "rpc"
)

type Generators []gen.Generator

type CommonConfig struct {
	Project project.Config
	Rpc     rpc.Config
}

func NewGenerator(config CommonConfig, typ Typ) (gen.Generator, error) {
	var (
		gn  gen.Generator
		err error
	)
	switch typ {
	case ProjectTyp:
		gn, err = project.NewProject(config.Project)
	case RpcType:
		gn, err = rpc.NewRpc(config.Rpc)
	default:
		err = fmt.Errorf("wrong generator type")
	}
	return gn, err
}

func NewGenerators(config CommonConfig, types []Typ) (Generators, error) {
	gens := make([]gen.Generator, 0, len(types))
	var (
		gn  gen.Generator
		err error
	)
	for _, typ := range types {
		if gn, err = NewGenerator(config, typ); err != nil {
			return nil, err
		}
		gens = append(gens, gn)
	}
	return gens, nil
}

func (g Generators) Generate(ctx context.Context) error {
	for _, gn := range g {
		if err := gn.Generate(ctx); err != nil {
			return err
		}
	}
	return nil
}
