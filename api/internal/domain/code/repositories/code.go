package repositories

import (
	"github.com/kodmain/thetiptop/api/internal/domain/code/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type CodeRepository struct{}

type CodeRepositoryInterface interface {
	ListErrors() map[string]*entities.Code
}

func (r *CodeRepository) ListErrors() map[string]*entities.Code {
	var codes = map[string]*entities.Code{}

	for key, value := range errors.ListErrors() {
		codes[key] = &entities.Code{
			Message: value.Error(),
			Code:    value.Code(),
		}
	}

	return codes
}

func NewCodeRepository() *CodeRepository {
	return &CodeRepository{}
}
