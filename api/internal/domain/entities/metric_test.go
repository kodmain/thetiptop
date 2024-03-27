package entities_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	method := "GET"
	path := "/api/v1/users"
	count := 10

	metrics := entities.NewMetric(method, path, count)

	assert.NotNil(t, metrics, "Metrics should not be nil")
	assert.Equal(t, method, metrics.Method, "Method should match")
	assert.Equal(t, []byte(path), metrics.Path, "Path should match")
	assert.Equal(t, count, metrics.Count, "Count should match")
}

func TestMetricMarshalProm(t *testing.T) {
	method := "GET"
	path := "/api/v1/users"
	count := 10

	metrics := entities.NewMetric(method, path, count)
	promMetrics := metrics.MarshalProm()

	assert.Len(t, promMetrics, 1)
	assert.Equal(t, method, promMetrics[0].Labels["method"], "Method should match")
	assert.Equal(t, path, promMetrics[0].Labels["path"], "Path should match")
	assert.Equal(t, float64(count), promMetrics[0].Value, "Value should match")
}
