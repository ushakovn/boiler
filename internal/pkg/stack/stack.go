package stack

import "sync"

type Stack[T any] struct {
  mtx   sync.Mutex
  elems []T
}

func NewStack[T any]() *Stack[T] {
  return &Stack[T]{}
}

func (s *Stack[T]) Pop() (T, bool) {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  if len(s.elems) == 0 {
    return *new(T), false
  }
  elem := s.elems[len(s.elems)-1]
  s.elems = s.elems[:len(s.elems)-1]

  return elem, true
}

func (s *Stack[T]) Push(elem T) {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  s.elems = append(s.elems, elem)
}

func (s *Stack[T]) Peek() (T, bool) {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  if len(s.elems) == 0 {
    return *new(T), false
  }
  return s.elems[len(s.elems)-1], true
}

func (s *Stack[T]) PeekWith(f func(T)) bool {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  if len(s.elems) == 0 {
    return false
  }
  f(s.elems[len(s.elems)-1])

  return true
}

func (s *Stack[T]) Elems() []T {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  return s.elems
}

func (s *Stack[T]) Elem(idx int) (T, bool) {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  if len(s.elems) <= idx {
    return *new(T), false
  }

  return s.elems[idx], true
}

func (s *Stack[T]) ElemWith(idx int, f func(T)) bool {
  s.mtx.Lock()
  defer s.mtx.Unlock()

  if len(s.elems) <= idx {
    return false
  }
  f(s.elems[idx])

  return true
}
