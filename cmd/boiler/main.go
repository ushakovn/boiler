package main

import (
	"context"
	"flag"
	"fmt"

	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ushakovn/boiler/internal/pkg/gens/factory"
)

func main() {
	logBoilerMark()
	ctx := context.Background()

	types, err := parseFlagTypes()
	if err != nil {
		log.Fatalf("boiler validation error: %v", err)
	}

	gens, err := factory.NewGenerators(factory.CommonConfig{}, types)
	if err != nil {
		log.Fatalf("boiler initialization error: %v", err)
	}

	log.Infof("boiler info: generation started")

	if err = gens.Generate(ctx); err != nil {
		log.Fatalf("boiler generation error: %v", err)
	}

	log.Infof("boiler info: generation finished")
}

func parseFlagTypes() ([]factory.Typ, error) {
	joinedTypes := flag.String("type", "", "comma separated generators types")
	flag.Parse()

	splitTypes := strings.Split(*joinedTypes, ",")
	countTypes := len(splitTypes)

	if countTypes == 0 || countTypes == 1 && splitTypes[0] == "" {
		return nil, fmt.Errorf("generator types not set: all")
	}
	types := make([]factory.Typ, 0, countTypes)

	for index, typ := range splitTypes {
		if typ == "" {
			return nil, fmt.Errorf("generator type not set: index=%d", index)
		}
		types = append(types, factory.Typ(typ))
	}
	return types, nil
}

func logBoilerMark() {
	log.Infof(`

 _           _ _           
| |         (_) |          
| |__   ___  _| | ___ _ __ 
| '_ \ / _ \| | |/ _ \ '__|
| |_) | (_) | | |  __/ |   
|_.__/ \___/|_|_|\___|_|


`)
}

type Request struct {
	Field string `json:"field"`
}

type CompositeType struct {
	Field string `json:"field"`
}
