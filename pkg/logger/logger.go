package logger

import (
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/env"
)

func SetDefaultLogOptions() {
  appEnv := env.AppEnv()

  if appEnv == env.LocalEnv {
    // Set text formatter for local run
    log.SetFormatter(&log.TextFormatter{
      TimestampFormat: "2006-01-02 15:04:05",
    })
    return
  }
  // Set JSON formatter otherwise
  log.SetFormatter(new(log.JSONFormatter))
}
