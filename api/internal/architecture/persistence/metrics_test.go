package persistence_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsRepository(t *testing.T) {
	var method = "GET"
	var path = "/api/v1/users"

	repo := persistence.NewMetricsRepository()
	assert.NotNil(t, repo, "Repository should not be nil")

	metrics, err := repo.GetAllRequestStats()
	assert.Error(t, err, "GetAllRequestStats should return an error")
	assert.Nil(t, metrics, "Metrics should be nil")
	assert.Len(t, metrics, 0, "Metrics length should be 1")

	metric, err := repo.GetMostFrequentRequest()
	assert.Error(t, err, "GetMostFrequentRequest should return an error")
	assert.Nil(t, metric, "Metric should not be nil")

	err = repo.IncrementRequestCount(method, path)
	assert.NoError(t, err, "IncrementRequestCount should not return an error")
	err = repo.IncrementRequestCount(method, path)
	assert.NoError(t, err, "IncrementRequestCount should not return an error")

	metrics, err = repo.GetAllRequestStats()
	assert.NoError(t, err, "GetAllRequestStats should not return an error")
	assert.NotNil(t, metrics, "Metrics should not be nil")
	assert.Len(t, metrics, 1, "Metrics length should be 1")

	metric, err = repo.GetMostFrequentRequest()
	assert.NoError(t, err, "GetMostFrequentRequest should not return an error")
	assert.NotNil(t, metric, "Metric should not be nil")
	assert.Equal(t, method, metric.Method, "Method should match")
	assert.Equal(t, path, string(metric.Path), "Path should match")
	assert.Equal(t, 2, metric.Count, "Count should match")

	assert.Equal(t, metrics[0], metric, "String representation should match")
}
