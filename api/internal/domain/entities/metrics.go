package entities

import "github.com/kodmain/thetiptop/api/internal/architecture/serializers/prom"

// Metrics represents the metrics for a specific request.
// It stores the HTTP method, the request path, and the count of requests.
type Metrics []*Metric

// NewMetric creates and returns a new Metrics instance.
// This function is a constructor for the Metrics entity, initializing it with provided values.
//
// Returns:
// - Metric: a new instance of Metrics.
func NewMetrics() Metrics {
	return make(Metrics, 0)
}

// MarshalProm convertit l'instance Metrics en MetricData.
//
// Returns:
// - MetricData: Une instance de MetricData représentant les métriques.
func (m Metrics) MarshalProm() []prom.MetricData {
	var metrics []prom.MetricData

	for _, metric := range m {
		metrics = append(metrics, metric.MarshalProm()...)
	}

	return metrics
}
