package httpgateway

import (
  "fmt"
  "net/http"

  "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func Handle(mux *runtime.ServeMux) func(method string, path string, handle http.HandlerFunc) error {
  return func(method string, path string, handler http.HandlerFunc) error {
    if err := mux.HandlePath(method, path, Wrap(handler)); err != nil {
      return fmt.Errorf("mux.HandlePath: %w", err)
    }
    return nil
  }
}

func Wrap(handler http.HandlerFunc) runtime.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
    handler(w, r)
  }
}
