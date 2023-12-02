package env

import "os"

const (
  ProductionEnv Env = "PRODUCTION"
  StagingEnv    Env = "STAGING"
  LocalEnv      Env = "LOCAL"
)

var knownAppEnvs = map[Env]struct{}{
  ProductionEnv: {},
  StagingEnv:    {},
  LocalEnv:      {},
}

type Env string

func AppEnv() Env {
  const envKey = "BOILER_APP_ENV"
  appEnv := Env(os.Getenv(envKey))

  if _, ok := knownAppEnvs[appEnv]; ok {
    return appEnv
  }
  return LocalEnv
}
