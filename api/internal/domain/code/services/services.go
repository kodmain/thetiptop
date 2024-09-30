package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/domain/code/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/code/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type CodeService struct {
	security security.PermissionInterface
	repo     repositories.CodeRepositoryInterface
}

func Code(security security.PermissionInterface, repo repositories.CodeRepositoryInterface) *CodeService {
	return &CodeService{security, repo}
}

type CodeServiceInterface interface {
	ListErrors() (map[string]*entities.Code, errors.ErrorInterface)
}
