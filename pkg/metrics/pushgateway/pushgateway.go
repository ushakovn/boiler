package pushgateway

import (
  "context"
  "fmt"
  "time"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/push"
  log "github.com/sirupsen/logrus"
)

type Pusher struct {
  duration   time.Duration
  prometheus *push.Pusher
}

type Config struct {
  Url      string
  Job      string
  Duration time.Duration
}

func NewPusher(config Config) *Pusher {
  return &Pusher{
    duration: config.Duration,

    prometheus: push.
      New(config.Url, config.Job).
      Gatherer(prometheus.DefaultGatherer),
  }
}

func (p *Pusher) Push(ctx context.Context) error {
  if err := p.prometheus.PushContext(ctx); err != nil {
    return fmt.Errorf("metrics push error: %w", err)
  }
  return nil
}

func (p *Pusher) Run(ctx context.Context) error {
  if err := p.prometheus.PushContext(ctx); err != nil {
    return fmt.Errorf("first metrics push error: %w", err)
  }

  go func() {
    ticker := time.NewTicker(p.duration)
    defer ticker.Stop()

    for {
      select {
      case <-ctx.Done():
        log.Infof("pushgateway.Run: pusher stopped: context cancelled")
        return

      case <-ticker.C:
        if err := p.prometheus.PushContext(ctx); err != nil {
          log.Errorf("pushgateway.Run: metrics push error: %v", err)
        } else {
          log.Debugf("pushgateway.Run: metrics push was sucessfull")
        }
      }
    }
  }()

  return nil
}
