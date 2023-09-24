package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/ushakovn/boiler/internal/pkg/gens/factory"
	"github.com/ushakovn/boiler/internal/pkg/gens/rpc"
)

func main() {
	// common generator flags
	genTypes := flag.String("types", "", "space separated generator type names")
	workDir := flag.String("dir", "", "working directory path")

	// project generator flags
	_ = flag.String("project", "./config/project/config.textproto", "project directory path")

	// controller generator flags
	apiDir := flag.String("api", "", "api contract directory path")

	// parse command line flags
	flag.Parse()

	parsedTypes := strings.Split(*genTypes, " ")
	parsedTyp := parsedTypes[0]

	gen, err := factory.NewGenerator(&factory.CommonConfig{
		Rpc: rpc.Config{
			WorkDir: *workDir,
			RpcDir:  *apiDir,
		},
	}, factory.Typ(parsedTyp))
	if err != nil {
		log.Fatalf("boiler initialize failed: %v", err)
	}

	if err := gen.Generate(context.Background()); err != nil {
		log.Fatalf("boiler generation failed: %v", err)
	}
}
