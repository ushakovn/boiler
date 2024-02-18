package gen

import (
  "github.com/spf13/cobra"
  cmdConfig "github.com/ushakovn/boiler/cmd/root/gen/config"
  cmdGqlgen "github.com/ushakovn/boiler/cmd/root/gen/gqlgen"
  cmdGrpc "github.com/ushakovn/boiler/cmd/root/gen/grpc"
  cmdKafkaoutbox "github.com/ushakovn/boiler/cmd/root/gen/kafkaoutbox"
  cmdClients "github.com/ushakovn/boiler/cmd/root/gen/protodeps"
  cmdStorage "github.com/ushakovn/boiler/cmd/root/gen/storage"
)

var CmdGen = &cobra.Command{
  Use: "gen",

  Short: "Generate components for a microservice application in the Go language",
  Long:  `Generate components for a microservice application in the Go language`,
}

func init() {
  CmdGen.AddCommand(
    cmdGrpc.CmdGrpc,
    cmdStorage.CmdStorage,
    cmdGqlgen.CmdGqlgen,
    cmdClients.CmdProtoDeps,
    cmdConfig.CmdConfig,
    cmdKafkaoutbox.CmdKafkaoutbox,
  )
}
