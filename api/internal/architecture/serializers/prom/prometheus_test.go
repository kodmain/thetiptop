package prom_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/architecture/serializers/prom"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

const (
	metricTest = "This is a test metric"
)

func TestNewMetricMeta(t *testing.T) {
	data := entities.NewMetrics()
	meta, err := prom.NewMetricMeta("test_metric", metricTest, "counter", data)
	assert.NoError(t, err, "NewMetricMeta should not return an error")
	assert.NotNil(t, meta, "MetricMeta should not be nil")
	assert.Equal(t, "test_metric", meta.Name, "Name should match")
	assert.Equal(t, metricTest, meta.Help, "Help should match")
	assert.Equal(t, "counter", meta.Type, "Type should match")
	assert.Len(t, meta.Value, 0, "Value length should be 0")
}

func TestMarshal(t *testing.T) {
	data := entities.NewMetrics()
	data = append(data, entities.NewMetric("GET", "/api/v1/users", 10))
	data = append(data, entities.NewMetric("POST", "/api/v1/users", 20))

	meta, err := prom.NewMetricMeta("test_metric", metricTest, "counter", data)
	assert.NoError(t, err, "NewMetricMeta should not return an error")
	assert.NotNil(t, meta, "MetricMeta should not be nil")

	b, err := prom.Marshal(meta)
	assert.NoError(t, err, "Marshal should not return an error")
	assert.NotNil(t, b, "Bytes should not be nil")
}

func TestNewSimpleMeta(t *testing.T) {
	name := "test_metric"
	help := metricTest
	metricType := "counter"
	value := 10
	labels := []prom.MetricLabel{
		{Key: "label1", Value: "value1"},
		{Key: "label2", Value: "value2"},
	}

	expectedMeta := &prom.MetricMeta{
		Name: name,
		Help: help,
		Type: metricType,
		Value: []prom.MetricData{
			{
				Labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
				Value: value,
			},
		},
	}

	meta := prom.NewSimpleMeta(name, help, metricType, value, labels...)

	assert.Equal(t, expectedMeta, meta, "MetricMeta should be equal to the expected value")
}
