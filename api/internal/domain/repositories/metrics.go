package repositories

import "github.com/kodmain/thetiptop/api/internal/domain/entities"

// MetricsRepository defines the interface for interacting with metrics storage.
type MetricsRepository interface {
	// IncrementRequestCount increments the count of requests for a given method and path.
	// Parameters:
	// - method: string - The HTTP method of the request (e.g., GET, POST).
	// - path: string - The path of the request.
	// Returns:
	// - error: possible error encountered while incrementing the request count.
	IncrementRequestCount(method, path string) error

	// GetMostFrequentRequest retrieves the metric for the most frequently made request.
	// Returns:
	// - *entities.Metric: The metric of the most frequent request, nil if not available.
	// - error: possible error encountered while retrieving the most frequent request.
	GetMostFrequentRequest() (*entities.Metric, error)

	// GetAllRequestStats retrieves statistics for all tracked requests.
	// Returns:
	// - entities.Metrics: A collection of metrics for all requests.
	// - error: possible error encountered while retrieving the request statistics.
	GetAllRequestStats() (entities.Metrics, error)
}
