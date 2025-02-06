package errors_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorsToErrorInterface(t *testing.T) {
	t.Parallel()

	errs := errors.Errors{}
	assert.Nil(t, errs.ToErrorInterface())

	err := errors.New(404, "not.found")
	errs["not_found"] = err
	assert.NotNil(t, errs.ToErrorInterface())
}

func TestErrorsError(t *testing.T) {
	t.Parallel()

	errs := errors.Errors{}
	assert.Equal(t, "{}", errs.Error())

	err := errors.New(404, "not.found")
	errs["not_found"] = err
	expected := `{"not_found":"not.found"}`
	assert.JSONEq(t, expected, errs.Error())
}

func TestErrorsCode(t *testing.T) {
	t.Parallel()

	errs := errors.Errors{}
	assert.Equal(t, 0, errs.Code())

	err := errors.New(404, "not.found")
	errs["not_found"] = err
	assert.Equal(t, 404, errs.Code())
}

func TestErrorsLog(t *testing.T) {
	t.Parallel()

	errs := errors.Errors{}
	err := errors.New(404, "not.found")
	errs["not_found"] = err

	assert.Nil(t, errs.Log(fmt.Errorf("error")))
}

func TestErrorsAdd(t *testing.T) {
	t.Parallel()

	errs := errors.Errors{}
	err := errors.New(404, "not.found")
	assert.NoError(t, errs.Add("not_found", err))

	duplicateErr := errors.New(500, "duplicate.error")
	assert.Error(t, errs.Add("not_found", duplicateErr))
}

func TestNewError(t *testing.T) {
	t.Parallel()

	err := errors.New(404, "not.found")
	assert.Equal(t, 404, err.Code())
	assert.Equal(t, "not.found", err.Error())

	errs := errors.ListErrors()
	assert.Equal(t, 42, len(errs))

	err.Log(fmt.Errorf("error"))
}

// TestMarshalJSON tests the MarshalJSON method of the Error struct.
func TestMarshalJSON(t *testing.T) {
	// Arrange: Create an instance of Error with sample data.
	err := errors.New(500, "Internal Server Error")

	// Act: Marshal the error into JSON.
	jsonData, jsonErr := json.Marshal(err)

	// Assert: Check for errors during marshalling.
	assert.NoError(t, jsonErr, "MarshalJSON should not return an error")

	// Convert JSON to string for comparison.
	actual := string(jsonData)

	// Define the expected JSON string.
	expected := `{"code":500,"message":"Internal Server Error"}`

	// Assert: Check if the actual JSON matches the expected JSON.
	assert.JSONEq(t, expected, actual, "JSON output mismatch")
}
