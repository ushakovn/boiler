package timer

import "time"

type Timer interface {
  Stop() bool
  Start() <-chan time.Time
}

type timer struct {
  timer    *time.Timer
  duration time.Duration
}

func NewTimer(duration time.Duration) Timer {
  return &timer{
    duration: duration,
    timer:    time.NewTimer(duration),
  }
}

func (t *timer) Stop() bool {
  return t.timer.Stop()
}

func (t *timer) Start() <-chan time.Time {
  defer t.timer.Reset(t.duration)
  return t.timer.C
}
