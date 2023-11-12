package aggr

func Filter[T any](t []T, f func(T) bool) []T {
  ts := make([]T, 0, len(t))

  for _, t := range t {
    if !f(t) {
      continue
    }
    ts = append(ts, t)
  }
  return ts
}
