package handler

import (
  "github.com/go-chi/chi/v5"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func WithMetricsHandler(router chi.Router) {
  router.Handle("/metrics", promhttp.Handler())
}
