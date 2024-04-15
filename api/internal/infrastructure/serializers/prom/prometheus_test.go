package prom_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/prom"
	"github.com/stretchr/testify/assert"
)

const (
	metricTest = "This is a test metric"
)

type Metric struct {
	Method string
	Path   []byte // The path is stored as a byte slice to avoid memory leaks.
	Count  int
}

type Metrics []*Metric

func (m *Metric) MarshalProm() []prom.MetricData {
	return []prom.MetricData{
		{
			Labels: map[string]string{"method": m.Method, "path": string(m.Path)},
			Value:  float64(m.Count),
		},
	}
}

func (m Metrics) MarshalProm() []prom.MetricData {
	var metrics []prom.MetricData

	for _, metric := range m {
		metrics = append(metrics, metric.MarshalProm()...)
	}

	return metrics
}

func TestNewMetricMeta(t *testing.T) {
	data := &Metric{}
	meta, err := prom.NewMetricMeta("test_metric", metricTest, "counter", data)
	assert.NoError(t, err, "NewMetricMeta should not return an error")
	assert.NotNil(t, meta, "MetricMeta should not be nil")
	assert.Equal(t, "test_metric", meta.Name, "Name should match")
	assert.Equal(t, metricTest, meta.Help, "Help should match")
	assert.Equal(t, "counter", meta.Type, "Type should match")
	assert.Len(t, meta.Value, 1, "Value length should be 0")
}

func TestMarshal(t *testing.T) {
	data := make(Metrics, 0)
	data = append(data, &Metric{Method: "GET", Path: []byte("/api/v1/users"), Count: 10})
	data = append(data, &Metric{Method: "POST", Path: []byte("/api/v1/users"), Count: 20})

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
