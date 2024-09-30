package services

import (
	"github.com/kodmain/thetiptop/api/internal/domain/code/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *CodeService) ListErrors() (map[string]*entities.Code, errors.ErrorInterface) {
	return s.repo.ListErrors(), nil
}
