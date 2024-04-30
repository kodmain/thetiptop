package data

import "testing"

func TestObjectGet(t *testing.T) {
	obj := Object{"key1": "value1", "key2": "value2"}

	// Test case 1: Existing key
	expected1 := "value1"
	if result1 := obj.Get("key1"); result1 != expected1 {
		t.Errorf("Expected %s, but got %s", expected1, result1)
	}

	// Test case 2: Non-existing key
	expected2 := ""
	if result2 := obj.Get("key3"); result2 != expected2 {
		t.Errorf("Expected %s, but got %s", expected2, result2)
	}
}
