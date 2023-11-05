package root

import (
  "os"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  cmdGen "github.com/ushakovn/boiler/cmd/root/gen"
  cmdInit "github.com/ushakovn/boiler/cmd/root/init"
)

var flagDebug bool

var CmdRoot = &cobra.Command{
  Short: "Boiler is a mini-framework for the development of microservices in the Go language",

  Long: `Boiler is a mini-framework for the development of microservices in the Go language
 _           _ _           
| |         (_) |          
| |__   ___  _| | ___ _ __ 
| '_ \ / _ \| | |/ _ \ '__|
| |_) | (_) | | |  __/ |   
|_.__/ \___/|_|_|\___|_|
`,

  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    log.SetFormatter(&log.TextFormatter{
      FullTimestamp:   true,
      TimestampFormat: "2006-01-02 15:04:05",
    })

    if flagDebug {
      log.SetLevel(log.DebugLevel)
    } else {
      log.SetLevel(log.InfoLevel)
    }
  },
}

func Execute() {
  const errExitCode = 1

  err := CmdRoot.Execute()
  if err != nil {
    os.Exit(errExitCode)
  }
}

func init() {
  CmdRoot.AddCommand(cmdInit.CmdInit, cmdGen.CmdGen)

  CmdRoot.PersistentFlags().BoolVar(&flagDebug, "debug", false, "sets debug logging level")
}
