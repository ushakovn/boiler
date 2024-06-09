package stack

import "sync"

type Stack[T any] struct {
  mu sync.Mutex
  el []T
}

func New[T any]() *Stack[T] {
  return new(Stack[T])
}

func (s *Stack[T]) Push(el T) {
  s.mu.Lock()
  defer s.mu.Unlock()

  s.el = append(s.el, el)
}

func (s *Stack[T]) Pop() T {
  s.mu.Lock()
  defer s.mu.Unlock()

  if len(s.el) == 0 {
    return *new(T)
  }

  el := s.el[len(s.el)-1]
  s.el = s.el[:len(s.el)-1]

  return el
}

func (s *Stack[T]) Peek() T {
  s.mu.Lock()
  defer s.mu.Unlock()

  if len(s.el) == 0 {
    return *new(T)
  }

  return s.el[len(s.el)-1]
}
