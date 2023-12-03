package metrics

import (
  "strings"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/ushakovn/boiler/pkg/config"
)

const (
  defaultNamespace = "boiler"
)

func subsystemFromConfig() string {
  name := config.ClientConfig().GetAppInfo().Name
  name = strings.ToLower(name)
  return name
}

func NewCounter(name, help string) prometheus.Counter {
  return promauto.NewCounter(prometheus.CounterOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
  })
}

func NewCounterVec(name, help string, labels []string) *prometheus.CounterVec {
  return promauto.NewCounterVec(prometheus.CounterOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
  },
    labels,
  )
}

func NewGauge(name, help string, labels []string) prometheus.Gauge {
  return promauto.NewGauge(prometheus.GaugeOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
  })
}

func NewGaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
  return promauto.NewGaugeVec(prometheus.GaugeOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
  },
    labels,
  )
}

func NewHistogram(name, help string, buckets []float64) prometheus.Histogram {
  return promauto.NewHistogram(prometheus.HistogramOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
    Buckets:   buckets,
  })
}

func NewHistogramVec(name, help string, buckets []float64, labels []string) *prometheus.HistogramVec {
  return promauto.NewHistogramVec(prometheus.HistogramOpts{
    Namespace: defaultNamespace,
    Subsystem: subsystemFromConfig(),
    Name:      name,
    Help:      help,
    Buckets:   buckets,
  },
    labels,
  )
}
