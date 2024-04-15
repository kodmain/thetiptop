package prom

import (
	"fmt"
	"strings"
)

// Marshalable is an interface for objects that can be marshaled into MetricData.
type Marshalable interface {
	// MarshalProm marshals the object into a slice of MetricData.
	MarshalProm() []MetricData
}

// MetricMeta holds metadata for a metric including its name, help text, type, and value.
type MetricMeta struct {
	Name  string       // Name of the metric.
	Help  string       // Help provides a description for the metric.
	Type  string       // Type of the metric (e.g., counter, gauge).
	Value []MetricData // Value contains the actual metric data to be marshaled.
}

// NewMetricMeta creates a new MetricMeta instance.
// Parameters:
// - name: string - The name of the metric.
// - help: string - Help text describing the metric.
// - metricType: string - The type of the metric (e.g., counter, gauge).
// - data: Marshalable - The object that can be marshaled into MetricData.
// Returns:
// - *MetricMeta: A pointer to the newly created MetricMeta.
// - error: An error that might occur during the creation of MetricMeta.
func NewMetricMeta(name, help, metricType string, data Marshalable) (*MetricMeta, error) {
	return &MetricMeta{
		Name:  name,
		Help:  help,
		Type:  metricType,
		Value: data.MarshalProm(),
	}, nil
}

// MetricData represents the data for a metric, consisting of labels and a value.
type MetricData struct {
	Labels map[string]string // Labels are key-value pairs that add context to the metric.
	Value  any               // Value is the actual data of the metric.
}

// MetricLabel represents the data for a metric, consisting of labels and a value.
type MetricLabel struct {
	Key   string // Key is the label key.
	Value string // Value is the label value.
}

// Marshal converts a slice of MetricMeta into a Prometheus-compatible byte slice.
// Parameters:
// - metas: []*MetricMeta - A slice of MetricMeta to be marshaled.
// Returns:
// - []byte: The marshaled data in Prometheus format.
// - error: An error that might occur during marshaling.
func Marshal(metas ...*MetricMeta) ([]byte, error) {
	var sb strings.Builder

	for _, meta := range metas {
		if meta != nil {
			sb.WriteString(fmt.Sprintf("# HELP %s %s\n", meta.Name, meta.Help))
			sb.WriteString(fmt.Sprintf("# TYPE %s %s\n", meta.Name, meta.Type))

			for _, value := range meta.Value {
				metricLine := meta.Name
				if len(value.Labels) > 0 {
					labels := make([]string, 0, len(value.Labels))
					for k, v := range value.Labels {
						labels = append(labels, fmt.Sprintf("%s=\"%s\"", k, v))
					}
					metricLine += fmt.Sprintf("{%s}", strings.Join(labels, ", "))
				}
				metricLine += fmt.Sprintf(" %v\n", value.Value)
				sb.WriteString(metricLine)
			}
		}
	}

	return []byte(sb.String()), nil
}

// NewSimpleMeta creates a new MetricData instance with a single value.
// Parameters:
// - name: string - The name of the metric.
// - help: string - Help text describing the metric.
// - metricType: string - The type of the metric (e.g., counter, gauge).
// - value: any - The value of the metric.
// - labels: map[string]string - Labels for the metric.
// Returns:
// - *MetricMeta: A pointer to the newly created MetricMeta.
func NewSimpleMeta(name, help, metricType string, value any, labels ...MetricLabel) *MetricMeta {
	labelMap := make(map[string]string, len(labels))
	for _, label := range labels {
		labelMap[label.Key] = label.Value
	}

	return &MetricMeta{
		Name: name,
		Help: help,
		Type: metricType,
		Value: []MetricData{
			{
				Labels: labelMap,
				Value:  value,
			},
		},
	}
}
