package services

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
)

func (s *UserService) UserAuth(dtoCredential *transfert.Credential) (*string, security.Role, errors.ErrorInterface) {
	if dtoCredential == nil {
		return nil, "", errors.ErrNoDto
	}

	// Lire les informations d'identification de l'utilisateur
	credential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: dtoCredential.Email,
	})

	if err != nil {
		return nil, "", err
	}

	// Comparer les hashs si les credentials existent
	if !credential.CompareHash(*dtoCredential.Password) {
		return nil, "", errors_domain_user.ErrCredentialNotValid
	}

	client, _, err := s.repo.ReadUser(&transfert.User{
		CredentialID: &credential.ID,
	})

	if err != nil {
		return nil, "", errors_domain_user.ErrUserNotFound
	}

	if client != nil {
		return &credential.ID, entities.ROLE_CLIENT, nil
	}

	return &credential.ID, entities.ROLE_EMPLOYEE, nil
}

func (s *UserService) PasswordUpdate(dto *transfert.Credential) errors.ErrorInterface {
	if dto == nil {
		return errors.ErrNoDto
	}

	credential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: dto.Email,
	})

	if err != nil {
		return err
	}

	password, err := hash.Hash(aws.String(*credential.Email+":"+*dto.Password), hash.BCRYPT)
	if err != nil {
		return err
	}

	credential.Password = password

	if err := s.repo.UpdateCredential(credential); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ValidationRecover(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) errors.ErrorInterface {
	if dtoValidation == nil || dtoCredential == nil {
		return errors.ErrNoDto
	}

	credential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: dtoCredential.Email,
	})

	if err != nil {
		return err
	}

	client, employee, err := s.repo.ReadUser(&transfert.User{
		CredentialID: &credential.ID,
	})

	if err != nil {
		return err
	}

	if client != nil {
		dtoValidation.ClientID = &client.ID
	}

	if employee != nil {
		dtoValidation.EmployeeID = &employee.ID
	}

	validation, err := s.repo.CreateValidation(dtoValidation)
	if err != nil {
		return err
	}

	if validation.Type != entities.PhoneValidation {
		go s.sendValidationMail(credential, validation)
	}

	return nil
}

// sendMail Send a templated email to a client
// This function handles the common logic for sending templated emails to clients.
//
// Parameters:
// - client: *entities.Client The client to send the email to.
// - templateName: string The name of the email template.
// - validationType: entities.ValidationType The type of validation to check.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) sendMail(credential *entities.Credential, validation *entities.Validation, templateName string) errors.ErrorInterface {
	tpl := template.NewTemplate(templateName)
	if tpl == nil {
		return errors.ErrMailTemplateNotFound
	}

	text, html, err := tpl.Inject(template.Data{
		"AppName": env.APP_NAME,
		"Token":   validation.Token.String(),
	})

	if err != nil {
		return err
	}

	subject := "The Tip Top"

	m := &mail.Mail{
		To:      []string{*credential.Email},
		Subject: subject,
		Text:    text,
		Html:    html,
	}

	for i := 0; i < 3; i++ {
		if err := s.mail.Send(m); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.ErrMailSendFailed
}

// sendValidationMail Send a signup confirmation email to a client
// This function sends a signup confirmation email to the specified client.
//
// Parameters:
// - client: *entities.Credential The client to send the email to.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) sendValidationMail(credential *entities.Credential, token *entities.Validation) errors.ErrorInterface {
	return s.sendMail(credential, token, "token")
}

// validateClientAndValidation Validate client and validation entities
// This function handles the common logic for validating client and validation entities.
//
// Parameters:
// - dtoValidation: *transfert.Validation The validation DTO.
// - dtoClient: *transfert.Client The client DTO.
//
// Returns:
// - validation: *entities.Validation The validated validation entity.
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) validateClientAndValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, errors.ErrorInterface) {
	if dtoValidation == nil || dtoCredential == nil {
		return nil, errors.ErrNoDto
	}

	validation, err := s.repo.ReadValidation(dtoValidation)
	if err != nil {
		return nil, err
	}

	if validation.ExpiresAt.Before(time.Now()) {
		return nil, errors_domain_user.ErrValidationExpired
	}

	if validation.Validated {
		return nil, errors_domain_user.ErrValidationAlreadyValidated
	}

	validation.Validated = true

	if err := s.repo.UpdateValidation(validation); err != nil {
		return nil, err
	}

	return validation, nil
}

// PasswordValidation Validate password recovery
// This function validates a password recovery request.
//
// Parameters:
// - dtoValidation: *transfert.Validation The validation DTO.
// - dtoClient: *transfert.Client The client DTO.
//
// Returns:
// - validation: *entities.Validation The validated validation entity.
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) PasswordValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, errors.ErrorInterface) {
	return s.validateClientAndValidation(dtoValidation, dtoCredential)
}

// MailValidation Validate sign-up
// This function validates a sign-up request.
//
// Parameters:
// - dtoValidation: *transfert.Validation The validation DTO.
// - dtoClient: *transfert.Client The client DTO.
//
// Returns:
// - validation: *entities.Validation The validated validation entity.
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) MailValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, errors.ErrorInterface) {
	return s.validateClientAndValidation(dtoValidation, dtoCredential)
}
