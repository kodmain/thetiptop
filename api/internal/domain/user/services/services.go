package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

type UserService struct {
	security security.PermissionInterface
	repo     repositories.UserRepositoryInterface
	mail     mail.ServiceInterface
}

func User(security security.PermissionInterface, repo repositories.UserRepositoryInterface, mail mail.ServiceInterface) *UserService {
	return &UserService{security, repo, mail}
}

type UserServiceInterface interface {
	// Credential
	UserAuth(dtoCredential *transfert.Credential) (*string, string, errors.ErrorInterface)
	PasswordUpdate(dtoCredential *transfert.Credential) errors.ErrorInterface
	ValidationRecover(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) errors.ErrorInterface
	PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, errors.ErrorInterface)
	MailValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, errors.ErrorInterface)

	// Client
	RegisterClient(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface)
	GetClient(dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface)
	DeleteClient(dtoClient *transfert.Client) errors.ErrorInterface
	UpdateClient(client *transfert.Client) (*entities.Client, errors.ErrorInterface)

	// Employee
	RegisterEmployee(dtoCredential *transfert.Credential, dtoEmployee *transfert.Employee) (*entities.Employee, errors.ErrorInterface)
	GetEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, errors.ErrorInterface)
	DeleteEmployee(dtoEmployee *transfert.Employee) errors.ErrorInterface
	UpdateEmployee(Employee *transfert.Employee) (*entities.Employee, errors.ErrorInterface)
}
