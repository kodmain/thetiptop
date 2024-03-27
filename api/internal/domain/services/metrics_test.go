package services_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsService(t *testing.T) {

	var method = "GET"
	var path = "/api/v1/users"

	var repository = persistence.NewMetricsRepository()
	var service = services.NewMetricsService(repository)

	assert.NotNil(t, service, "Service should not be nil")

	err := service.IncrementRequest(method, path)
	assert.NoError(t, err, "IncrementRequest should not return an error")

	metrics, err := service.GetAllRequestStats()
	assert.NoError(t, err, "GetAllRequestStats should not return an error")
	assert.NotNil(t, metrics, "Metrics should not be nil")
	assert.Len(t, metrics, 1, "Metrics length should be 1")

	metric, err := service.GetMostFrequentRequest()
	assert.NoError(t, err, "GetMostFrequentRequest should not return an error")
	assert.NotNil(t, metric, "Metric should not be nil")
	assert.Equal(t, method, metric.Method, "Method should match")
	assert.Equal(t, path, string(metric.Path), "Path should match")
	assert.Equal(t, 1, metric.Count, "Count should match")

	assert.Equal(t, metrics[0], metric, "String representation should match")
}
