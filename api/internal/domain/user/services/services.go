package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

type UserService struct {
	repo repositories.UserRepositoryInterface
	mail mail.ServiceInterface
}

func User(repo repositories.UserRepositoryInterface, mail mail.ServiceInterface) *UserService {
	return &UserService{repo, mail}
}

type UserServiceInterface interface {
	// Credential
	UserAuth(dtoCredential *transfert.Credential) (*string, error)
	PasswordUpdate(dtoCredential *transfert.Credential) error
	ValidationRecover(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) error
	PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)
	MailValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)

	// Client
	RegisterClient(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, error)
	GetClient(dtoClient *transfert.Client) (*entities.Client, error)
	DeleteClient(dtoClient *transfert.Client) error
	UpdateClient(client *transfert.Client) (*entities.Client, error)

	// Employee
	RegisterEmployee(dtoCredential *transfert.Credential, dtoEmployee *transfert.Employee) (*entities.Employee, error)
	GetEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, error)
	DeleteEmployee(dtoEmployee *transfert.Employee) error
	UpdateEmployee(Employee *transfert.Employee) (*entities.Employee, error)
}
