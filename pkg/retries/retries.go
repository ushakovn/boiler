package retries

import (
  "context"
  "fmt"
  "time"

  log "github.com/sirupsen/logrus"
  "github.com/ushakovn/boiler/pkg/timer"
)

type Options struct {
  Count int
  Wait  time.Duration
}

func (o Options) Validate() error {
  if o.Count <= 0 {
    return fmt.Errorf("count must be positive")
  }
  if o.Wait <= 0 {
    return fmt.Errorf("wait must be positive")
  }
  return nil
}

func (o Options) WithDefault() Options {
  const (
    count = 5
    wait  = 100 * time.Millisecond
  )
  if o.Count == 0 {
    o.Count = count
  }
  if o.Wait == 0 {
    o.Wait = wait
  }
  return o
}

func DoWithRetries(ctx context.Context, opt Options, f func(ctx context.Context) error) error {
  opt = opt.WithDefault()

  tm := timer.NewTimer(opt.Wait)
  defer tm.Stop()

  var (
    index int
    err   error
  )

  for index < opt.Count {
    select {
    case <-ctx.Done():
      return fmt.Errorf("context cancelled")

    case <-tm.Start():
      if err = f(ctx); err != nil {
        log.Errorf("DoWithRetries: %v", err)
        index++
        continue
      }
      return nil
    }
  }

  return err
}
