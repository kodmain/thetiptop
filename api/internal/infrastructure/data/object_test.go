package data_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestObjectGet(t *testing.T) {
	v1 := "value1"
	v2 := "value2"

	obj := data.Object{
		"key1": &v1,
		"key2": &v2,
	}

	if result1 := obj.Get("key1"); result1 != nil {
		if str, ok := result1.(*string); ok {
			if *str != v1 {
				t.Errorf("expected %s, but got %s", v1, *str)
			}
		} else {
			t.Errorf("expected type *string, but got %T", result1)
		}
	}

	if result2 := obj.Get("key2"); result2 != nil {
		if str, ok := result2.(*string); ok {
			if *str != v2 {
				t.Errorf("expected %s, but got %s", v2, *str)
			}
		} else {
			t.Errorf("expected type *string, but got %T", result2)
		}
	}

	assert.Nil(t, obj.Get("key3"))
}

func TestObject_Hydrate(t *testing.T) {
	type Target struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}

	tests := []struct {
		name    string
		object  data.Object
		want    Target
		wantErr bool
	}{
		{
			name: "successful hydration",
			object: data.Object{
				"key1": newString("value1"),
				"key2": newString("value2"),
			},
			want: Target{
				Key1: "value1",
				Key2: "value2",
			},
			wantErr: false,
		},
		{
			name: "partial hydration",
			object: data.Object{
				"key1": newString("value1"),
			},
			want: Target{
				Key1: "value1",
				Key2: "",
			},
			wantErr: false,
		},
		{
			name:   "empty object",
			object: data.Object{},
			want: Target{
				Key1: "",
				Key2: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target Target
			err := tt.object.Hydrate(&target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hydrate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
