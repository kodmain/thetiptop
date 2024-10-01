package errors_test

import (
	"fmt"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Parallel()

	err := errors.New(404, "not.found")
	assert.Equal(t, 404, err.Code())
	assert.Equal(t, "not.found", err.Error())

	errs := errors.ListErrors()
	assert.Equal(t, 39, len(errs))

	err.Log(fmt.Errorf("error"))
}
