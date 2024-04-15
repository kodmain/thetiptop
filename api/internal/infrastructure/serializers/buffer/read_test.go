package buffer_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/buffer"
)

// TestReadFunction tests the Read function.
//
// Parameters:
// - t: *testing.T - Testing context provided by Go testing framework.
func TestReadFunction(t *testing.T) {
	// Create a test input buffer
	input := bytes.NewBufferString("test data")

	// Create a mock implementation of io.ReadCloser
	mockReader := io.NopCloser(input)

	// Call the Read function with the mock reader
	result, err := buffer.Read(mockReader)

	// Check if there's an error returned
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Check if the result matches the expected value
	expected := "test data"
	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result.String())
	}
}

func TestReadFunctionError(t *testing.T) {
	// Create a mock implementation of io.ReadCloser that returns an error when copied from
	mockReader := &mockReaderWithError{}

	// Call the Read function with the mock reader
	_, err := buffer.Read(mockReader)

	// Check if there's an error returned
	if err == nil {
		t.Error("Expected an error, but got nil")
		return
	}

	// Test passed
	t.Logf("Error encountered: %v", err)
}

// mockReaderWithError is a mock implementation of io.ReadCloser that always returns an error when copied from.
type mockReaderWithError struct{}

// Read method returns an error.
func (m *mockReaderWithError) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error: unable to read from reader")
}

// Close method does nothing.
func (m *mockReaderWithError) Close() error {
	return nil
}
