package services

import (
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/repositories"
)

// MetricsService is a service that provides operations for tracking and retrieving
// request metrics.
type MetricsService struct {
	// repo is the repository responsible for managing metrics data.
	repo repositories.MetricsRepository
}

// NewMetricsService creates and returns a new instance of MetricsService.
// It requires a MetricsRepository implementation to interact with the metrics data.
//
// Parameters:
// - repo: MetricsRepository - A repository that provides data access functionalities for metrics.
//
// Returns:
// - *MetricsService: A pointer to the newly created MetricsService instance.
func NewMetricsService(repo repositories.MetricsRepository) *MetricsService {
	return &MetricsService{
		repo: repo,
	}
}

// IncrementRequest increments the count for a given request method and path.
// It acts as a facade to the repository's IncrementRequestCount function,
// simplifying its usage in the application.
//
// Parameters:
// - method: string - The HTTP method of the request (e.g., GET, POST).
// - path: string - The path of the request.
//
// Returns:
// - error: An error if the increment operation fails, otherwise nil.
func (s *MetricsService) IncrementRequest(method, path string) error {
	return s.repo.IncrementRequestCount(method, path)
}

// GetMostFrequentRequest retrieves the metrics data for the most frequently made request.
// It delegates to the MetricsRepository to fetch this data and returns it to the caller.
//
// Returns:
// - *entities.Metric: The metrics data of the most frequent request.
// - error: An error if the retrieval fails, otherwise nil.
func (s *MetricsService) GetMostFrequentRequest() (*entities.Metric, error) {
	return s.repo.GetMostFrequentRequest()
}

// GetAllRequestStats retrieves the statistics for all tracked requests.
// It uses the MetricsRepository to fetch this data and returns it to the caller.
//
// Returns:
// - entities.Metrics: A collection of metrics for all requests.
// - error: An error if the retrieval fails, otherwise nil.
func (s *MetricsService) GetAllRequestStats() (entities.Metrics, error) {
	return s.repo.GetAllRequestStats()
}
