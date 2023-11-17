package gqlgen

import (
  "net/http"
  "text/template"

  "github.com/ushakovn/boiler/internal/pkg/templater"
  "github.com/ushakovn/boiler/templates"
)

func SandboxHandler(title string, endpoint string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/html; charset=UTF-8")
    var dummy template.FuncMap

    sandbox := templater.MustExecTemplate(templates.GqlgenSandbox, map[string]any{
      "Title":           title,
      "InitialEndpoint": endpoint,
    }, dummy)
    if _, err := w.Write(sandbox); err != nil {
      panic(err)
    }
  }
}
