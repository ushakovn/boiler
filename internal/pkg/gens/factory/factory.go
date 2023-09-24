package factory

import (
	"fmt"

	"github.com/ushakovn/boiler/internal/app/gen"
	"github.com/ushakovn/boiler/internal/pkg/gens/project"
	"github.com/ushakovn/boiler/internal/pkg/gens/rpc"
)

type Typ string

const (
	ProjectTyp Typ = "project"
	RpcType    Typ = "rpc"
)

type CommonConfig struct {
	Project project.Config
	Rpc     rpc.Config
}

func NewGenerator(config *CommonConfig, typ Typ) (gen.Generator, error) {
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
