package pushgateway

import (
  "context"
  "fmt"
  "time"

  validation "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
  "github.com/prometheus/client_golang/prometheus/push"
  log "github.com/sirupsen/logrus"
)

type RunParams struct {
  Url      string
  Job      string
  Duration time.Duration
}

func (p RunParams) Validate() error {
  return validation.ValidateStruct(&p,
    validation.Field(p.Url, is.URL),
    validation.Field(p.Job, validation.Required),
    validation.Field(p.Duration, validation.Required),
  )
}

func Run(ctx context.Context, params RunParams) error {
  pusher := push.New(params.Url, params.Job)

  if err := pusher.PushContext(ctx); err != nil {
    return fmt.Errorf("first metrics push error: %w", err)
  }

  go func() {
    ticker := time.NewTicker(params.Duration)
    defer ticker.Stop()

    for {
      select {
      case <-ctx.Done():
        if err := pusher.PushContext(ctx); err != nil {
          log.Errorf("pushgateway.Run: last metrics push error: %v", err)
        } else {
          log.Infof("pushgateway.Run: last metrics push was sucessfull")
        }

        log.Warnf("pushgateway.Run: pusher stopped: context cancelled")
        return

      case <-ticker.C:
        if err := pusher.PushContext(ctx); err != nil {
          log.Errorf("pushgateway.Run: ticker metrics push error: %v", err)
        } else {
          log.Infof("pushgateway.Run: ticker metrics push was sucessfull")
        }
      }
    }
  }()

  return nil
}
