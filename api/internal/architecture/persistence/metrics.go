package persistence

import (
	"errors"
	"sync"

	"github.com/kodmain/thetiptop/api/internal/domain/entities"
)

// MetricsRepository is a repository for storing and retrieving metrics data.
type MetricsRepository struct {
	mu             sync.Mutex                  // Mutex to ensure thread-safe access to the requestCounter.
	requestCounter map[string]*entities.Metric // requestCounter holds the count of requests by method and path.
}

// NewMetricsRepository creates and returns a new instance of MetricsRepository.
func NewMetricsRepository() *MetricsRepository {
	return &MetricsRepository{
		requestCounter: make(map[string]*entities.Metric),
	}
}

// IncrementRequestCount increases the count for a specific request method and path.
// Parameters:
// - method: string - The HTTP method of the request (e.g., GET, POST).
// - path: string - The path of the request.
// Returns:
// - error: Possible error encountered while incrementing the request count.
func (r *MetricsRepository) IncrementRequestCount(method, path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := method + ":" + path
	stats, exists := r.requestCounter[key]

	if exists {
		stats.Count++
	} else {
		r.requestCounter[key] = entities.NewMetric(method, path, 1)
	}

	return nil
}

// GetMostFrequentRequest retrieves the metric for the most frequently made request.
// Returns:
// - *entities.Metric: The metric for the most frequent request or nil if no data is available.
// - error: Possible error encountered while retrieving the metric.
func (r *MetricsRepository) GetMostFrequentRequest() (*entities.Metric, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var mostFrequent *entities.Metric

	for _, stats := range r.requestCounter {
		if mostFrequent == nil || stats.Count > mostFrequent.Count {
			mostFrequent = stats
		}
	}

	if mostFrequent == nil {
		return nil, errors.New("no request has been made yet")
	}

	return mostFrequent, nil
}

// GetAllRequestStats retrieves all request statistics.
// Returns:
// - entities.Metrics: A slice of metrics representing all tracked requests.
// - error: Possible error encountered while retrieving the statistics.
func (r *MetricsRepository) GetAllRequestStats() (entities.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.requestCounter) == 0 {
		return nil, errors.New("no request has been made yet")
	}

	allStats := make(entities.Metrics, 0, len(r.requestCounter))

	for _, stats := range r.requestCounter {
		allStats = append(allStats, stats)
	}

	return allStats, nil
}
