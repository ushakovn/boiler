package closer

import (
  "os"
  "os/signal"
  "sync"

  log "github.com/sirupsen/logrus"
)

type Closer interface {
  Add(f func() error)
  WaitAll()
  CloseAll()
}

type closer struct {
  mu    sync.Mutex
  once  sync.Once
  calls []func() error
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

func (c *closer) Add(f func() error) {
  c.mu.Lock()
  defer c.mu.Unlock()
  c.calls = append(c.calls, f)
}

func (c *closer) WaitAll() {
  <-c.done
}

func (c *closer) CloseAll() {
  c.once.Do(func() {
    c.mu.Lock()
    defer c.mu.Unlock()

    callsCount := len(c.calls)
    errCh := make(chan error, callsCount)

    for _, call := range c.calls {
      go func(c func() error) {
        errCh <- c()
      }(call)
    }
    var errCount int

    for i := 0; i < callsCount; i++ {
      if err := <-errCh; err != nil {
        log.Warnf("boiler: closer call error: %v", err)
        errCount++
      }
    }

    if errCount == 0 {
      log.Info("boiler: closer close all finished with success")
    } else {
      log.Warnf("boiler: closer close all finished with %d errors", errCount)
    }

    c.done <- struct{}{}
  })
}
