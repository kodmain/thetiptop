package repositories_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/code/repositories"
	"github.com/stretchr/testify/assert"
)

func TestCodeRepository_ListErrors(t *testing.T) {
	repo := repositories.NewCodeRepository()
	codes := repo.ListErrors()
	assert.NotNil(t, codes)
}
