package entities

import "github.com/kodmain/thetiptop/api/internal/architecture/serializers/prom"

// Metrics represents the metrics for a specific request.
// It stores the HTTP method, the request path, and the count of requests.
type Metric struct {
	Method string
	Path   []byte // The path is stored as a byte slice to avoid memory leaks.
	Count  int
}

// NewMetric creates and returns a new Metrics instance.
// This function is a constructor for the Metrics entity, initializing it with provided values.
//
// Parameters:
// - method: string - The HTTP method of the request (e.g., GET, POST).
// - path: string - The path of the request.
// - count: int - The initial count of requests for the given method and path.
//
// Returns:
// - *Metrics: A pointer to the newly created Metrics instance.
func NewMetric(method, path string, count int) *Metric {
	return &Metric{
		Method: method,
		Path:   []byte(path),
		Count:  count,
	}
}

// MarshalProm convertit l'instance Metrics en MetricData.
//
// Returns:
// - MetricData: Une instance de MetricData représentant les métriques.
func (m *Metric) MarshalProm() []prom.MetricData {
	return []prom.MetricData{
		{
			Labels: map[string]string{"method": m.Method, "path": string(m.Path)},
			Value:  float64(m.Count),
		},
	}
}
