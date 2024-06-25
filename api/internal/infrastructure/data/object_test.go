package data_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestObjectGet(t *testing.T) {
	obj := data.Object{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}

	// Test case 1: Existing key
	expected1 := aws.String("value1")
	if result1 := obj.Get("key1"); *result1 != *expected1 {
		t.Errorf("expected %s, but got %s", *expected1, *result1)
	}

	// Test case 2: Non-existing key
	assert.Nil(t, obj.Get("key3"))
}
