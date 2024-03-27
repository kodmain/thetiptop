package entities_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	metrics := entities.NewMetrics()
	assert.NotNil(t, metrics, "Metrics should not be nil")
	assert.Len(t, metrics, 0)
}

func TestMetricsMarshalProm(t *testing.T) {
	metrics := entities.NewMetrics()

	metrics = append(metrics, entities.NewMetric("GET", "/api/v1/users", 10))
	metrics = append(metrics, entities.NewMetric("POST", "/api/v1/users", 20))

	promMetrics := metrics.MarshalProm()
	assert.Len(t, promMetrics, 2)
}
