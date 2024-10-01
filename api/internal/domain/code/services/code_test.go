package services_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/code/entities"
	"github.com/stretchr/testify/assert"
)

func TestCodeService_ListErrors(t *testing.T) {
	service, _, _, mockRepository := setup()

	mockErrorMap := map[string]*entities.Code{
		"Error1": {
			Code:    400,
			Message: "Bad Request",
		},
		"Error2": {
			Code:    404,
			Message: "Not Found",
		},
	}

	mockRepository.On("ListErrors").Return(mockErrorMap, nil)

	codes, err := service.ListErrors()

	assert.Nil(t, err)
	assert.NotNil(t, codes)
	assert.Equal(t, 2, len(codes))

	assert.Equal(t, 400, codes["Error1"].Code)
	assert.Equal(t, "Bad Request", codes["Error1"].Message)
	assert.Equal(t, 404, codes["Error2"].Code)
	assert.Equal(t, "Not Found", codes["Error2"].Message)

	mockRepository.AssertExpectations(t)
}
