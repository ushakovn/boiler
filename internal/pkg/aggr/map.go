package aggr

func ToMap[K comparable, V any](t []V, f func(V) K) map[K]V {
  m := make(map[K]V)

  for _, t := range t {
    k := f(t)
    m[k] = t
  }
  return m
}

func Map[T any, S any](t []T, f func(T) S) []S {
  s := make([]S, 0, len(t))

  for _, t := range t {
    s = append(s, f(t))
  }
  return s
}

func MapLookup[K comparable, V any, M map[K]V](m M, k K) bool {
  _, ok := m[k]
  return ok
}
