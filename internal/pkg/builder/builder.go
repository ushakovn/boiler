package builder

import (
  "fmt"
  "strings"
  "sync"
)

type Builder struct {
  mu sync.Mutex
  s  []string
}

func NewBuilder() Builder {
  return Builder{}
}

func (b *Builder) Clear() {
  b.mu.Lock()
  defer b.mu.Unlock()
  b.s = nil
}

func (b *Builder) Write(f string, s ...any) {
  b.mu.Lock()
  defer b.mu.Unlock()

  f = fmt.Sprintf(f, s...)
  b.s = append(b.s, fmt.Sprintf("%v", f))
}

func (b *Builder) Count() int {
  b.mu.Lock()
  defer b.mu.Unlock()

  return len(b.s)
}

func (b *Builder) String() string {
  b.mu.Lock()
  defer b.mu.Unlock()
  return strings.Join(b.s, "")
}
