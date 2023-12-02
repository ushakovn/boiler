package closer

import (
  "context"
  "math"
  "os"
  "os/signal"
  "sync"

  log "github.com/sirupsen/logrus"
)

type Closer interface {
  Add(f ...func(context.Context) error)
  WaitAll()
  CloseAll()
}

type closer struct {
  mu    sync.Mutex
  once  sync.Once
  calls []func(context.Context) error
  done  chan struct{}
}

func NewCloser(signals ...os.Signal) Closer {
  c := &closer{
    done: make(chan struct{}),
  }
  if len(signals) > 0 {
    go func(c *closer) {
      ch := make(chan os.Signal, 1)
      defer close(ch)

      signal.Notify(ch, signals...)
      <-ch
      signal.Stop(ch)

      c.CloseAll()
    }(c)
  }
  return c
}

func (c *closer) Add(f ...func(context.Context) error) {
  c.mu.Lock()
  defer c.mu.Unlock()

  c.calls = append(c.calls, f...)
}

func (c *closer) WaitAll() {
  <-c.done
}

func (c *closer) CloseAll() {
  const maxWatchers = 10
  ctx := context.Background()

  c.once.Do(func() {
    c.mu.Lock()
    defer c.mu.Unlock()

    errCh := make(chan error)
    var errCount int

    callsCount := len(c.calls)
    watchersCount := int(math.Min(float64(callsCount), maxWatchers))

    for i := 0; i < watchersCount; i++ {
      go func() {
        for err := range errCh {
          if err == nil {
            continue
          }
          log.Warnf("boiler: closer call error: %v", err)
          errCount++
        }
      }()
    }

    wg := sync.WaitGroup{}

    for _, call := range c.calls {
      wg.Add(1)

      go func(c func(ctx context.Context) error) {
        defer wg.Done()
        errCh <- c(ctx)

      }(call)
    }

    wg.Wait()
    close(errCh)

    if errCount == 0 {
      log.Info("boiler: closer close all with success")
    } else {
      log.Warnf("boiler: closer close all with %d errors", errCount)
    }

    c.done <- struct{}{}
  })
}
