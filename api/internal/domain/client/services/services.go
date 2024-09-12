package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
)

type ClientServiceInterface interface {
	// User
	UserRegister(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, error)
	UserAuth(dtoCredential *transfert.Credential) (*entities.Client, error)
	PasswordUpdate(dtoCredential *transfert.Credential) error
	ValidationRecover(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) error

	// Validation
	PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)
	MailValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)

	// Client
	GetClient(dtoClient *transfert.Client) (*entities.Client, error)
	UpdateClient(client *transfert.Client) error
}
