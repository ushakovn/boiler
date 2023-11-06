package gen

import (
  "github.com/spf13/cobra"
  cmdGqlgen "github.com/ushakovn/boiler/cmd/root/gen/gqlgen"
  cmdGrpc "github.com/ushakovn/boiler/cmd/root/gen/grpc"
  cmdStorage "github.com/ushakovn/boiler/cmd/root/gen/storage"
)

var CmdGen = &cobra.Command{
  Use: "gen",

  SuggestFor: []string{
    "generate",
  },

  Short: "Generate components for a microservice application in the Go language",
  Long:  `Generate components for a microservice application in the Go language`,
}

func init() {
  CmdGen.AddCommand(cmdGrpc.CmdGrpc, cmdStorage.CmdStorage, cmdGqlgen.CmdGqlgen)
}
