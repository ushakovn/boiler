package env

const (
  AppEnvKey Key = "BOILER_APP_ENV"
)

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

func AppEnv() Env {
  env := Get(AppEnvKey)

  if _, ok := knownAppEnvs[env]; ok {
    return env
  }
  return LocalEnv
}
