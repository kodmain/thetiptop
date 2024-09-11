package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
)

type ClientServiceInterface interface {
	UserRegister(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, error)
	UserAuth(dtoCredential *transfert.Credential) (*entities.Client, error)
	PasswordUpdate(dtoCredential *transfert.Credential) error
	ValidationRecover(dtoValidation *transfert.Validation, obj *transfert.Credential) error
	SignValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)
	PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Credential) (*entities.Validation, error)

	UpdateClient(client *transfert.Client) error
}
