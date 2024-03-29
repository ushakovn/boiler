package init

import (
  "github.com/spf13/cobra"
  cmdConfig "github.com/ushakovn/boiler/cmd/root/init/config"
  cmdGqlgen "github.com/ushakovn/boiler/cmd/root/init/gqlgen"
  cmdGrpc "github.com/ushakovn/boiler/cmd/root/init/grpc"
  cmdKafkaoutbox "github.com/ushakovn/boiler/cmd/root/init/kafkaoutbox"
  cmdProtoDeps "github.com/ushakovn/boiler/cmd/root/init/protodeps"
  cmdStorage "github.com/ushakovn/boiler/cmd/root/init/storage"
)

var CmdInit = &cobra.Command{
  Use: "init",

  Short: "Init a template for a microservice application in the Go language",
  Long:  `Init a template for a microservice application in the Go language`,
}

func init() {
  CmdInit.AddCommand(
    cmdGrpc.CmdGrpc,
    cmdGqlgen.CmdGqlgen,
    cmdProtoDeps.CmdProtoDeps,
    cmdStorage.CmdStorage,
    cmdConfig.CmdConfig,
    cmdKafkaoutbox.CmdKafkaoutbox,
  )
}
