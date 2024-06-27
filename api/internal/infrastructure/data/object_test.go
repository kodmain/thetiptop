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

	if result1 := obj.Get("key1"); *result1 != v1 {
		t.Errorf("expected %s, but got %s", v1, *result1)
	}

	assert.Nil(t, obj.Get("key3"))
}
