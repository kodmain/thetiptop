package database_test

import (
	"reflect"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string  `json:"field1"`
	Field2 int     `json:"field2"`
	Field3 *string `json:"field3"`
	Field4 *int    `json:"field4"`
}

func TestStructToMap(t *testing.T) {
	field3Value := "value3"
	field4Value := 4

	tests := []struct {
		name   string
		input  any
		output map[string]any
	}{
		{
			name: "All fields set",
			input: TestStruct{
				Field1: "value1",
				Field2: 2,
				Field3: &field3Value,
				Field4: &field4Value,
			},
			output: map[string]any{
				"field1": "value1",
				"field2": 2,
				"field3": "value3",
				"field4": 4,
			},
		},
		{
			name: "Nil pointer fields",
			input: TestStruct{
				Field1: "value1",
				Field2: 2,
				Field3: nil,
				Field4: nil,
			},
			output: map[string]any{
				"field1": "value1",
				"field2": 2,
				"field3": nil,
				"field4": nil,
			},
		},
		{
			name:  "Empty struct",
			input: TestStruct{},
			output: map[string]any{
				"field1": "",
				"field2": 0,
				"field3": nil,
				"field4": nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := database.StructToMap(tt.input)
			assert.True(t, reflect.DeepEqual(result, tt.output))
		})
	}
}
