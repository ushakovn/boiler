package root

import (
  "context"
  "fmt"
  "os"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  cmdGen "github.com/ushakovn/boiler/cmd/root/gen"
  cmdInit "github.com/ushakovn/boiler/cmd/root/init"
  "github.com/ushakovn/boiler/internal/pkg/executor"
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

  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

    return execModTidy()
  },

  PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
    return execModTidy()
  },
}

func Execute() {
  const errExitCode = 1
  setLogFormatter()

  if err := CmdRoot.Execute(); err != nil {
    os.Exit(errExitCode)
  }
}

func init() {
  CmdRoot.AddCommand(cmdInit.CmdInit, cmdGen.CmdGen)

  CmdRoot.PersistentFlags().BoolVar(&flagDebug, "debug", false, "sets debug logging level")
}

func setLogFormatter() {
  log.SetFormatter(&log.TextFormatter{
    FullTimestamp:   true,
    TimestampFormat: "2006-01-02 15:04:05",
  })
  if flagDebug {
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }
}

func execModTidy() error {
  if err := executor.ExecCommandContext(context.Background(), "go", "mod", "tidy"); err != nil {
    return fmt.Errorf("boiler: failed to exec go mod tidy: %w", err)
  }
  return nil
}