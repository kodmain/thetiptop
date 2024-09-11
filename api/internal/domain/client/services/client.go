package services

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
)

type ClientService struct {
	repo repositories.ClientRepositoryInterface
	mail mail.ServiceInterface
}

func Client(repo repositories.ClientRepositoryInterface, mail mail.ServiceInterface) *ClientService {
	return &ClientService{repo, mail}
}

func (s *ClientService) UserRegister(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, error) {
	if dtoCredential == nil || dtoClient == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	logger.Info("UserRegister", "dtoCredential", dtoCredential, "dtoClient", dtoClient)
	_, err := s.repo.ReadCredential(dtoCredential)
	if err == nil {
		return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
	}

	credential, err := s.repo.CreateCredential(dtoCredential)
	if err != nil {
		return nil, err
	}

	client, err := s.repo.CreateClient(dtoClient)
	if err != nil {
		return nil, err
	}

	credential.ClientID = &client.ID

	client.Validations = append(client.Validations, &entities.Validation{
		ClientID: &client.ID,
		Type:     entities.MailValidation,
	})

	if err := s.repo.UpdateClient(client); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateCredential(credential); err != nil {
		return nil, err
	}

	go s.sendValidationMail(credential, client.Validations[0])

	return client, nil
}

func (s *ClientService) UserAuth(credential *transfert.Credential) (*entities.Client, error) {
	if credential == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	// Lire les informations d'identification de l'utilisateur
	clientCredential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: credential.Email,
	})

	if err != nil {
		return nil, err // Erreur si les credentials ne sont pas trouvés
	}

	// Comparer les hashs si les credentials existent
	if !clientCredential.CompareHash(*credential.Password) {
		return nil, fmt.Errorf(errors.ErrCredentialNotFound)
	}

	client, err := s.repo.ReadClient(&transfert.Client{
		ID: clientCredential.ClientID,
	})

	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	// Continuer le processus d'authentification
	return client, nil
}

func (s *ClientService) UpdateClient(dtoClient *transfert.Client) error {
	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
	}

	if err := s.repo.UpdateClient(client); err != nil {
		return err
	}

	return nil
}

func (s *ClientService) ValidationRecover(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) error {
	if dtoValidation == nil || dtoCredential == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	credential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: dtoCredential.Email,
	})

	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
	}

	client, err := s.repo.ReadClient(&transfert.Client{
		ID: credential.ClientID,
	})

	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
	}

	dtoValidation.ClientID = &client.ID
	validation, err := s.repo.CreateValidation(dtoValidation)
	if err != nil {
		return err
	}

	if validation.Type != entities.PhoneValidation {
		go s.sendValidationMail(credential, validation)
	}

	return nil
}

func (s *ClientService) PasswordUpdate(dto *transfert.Credential) error {
	if dto == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	credential, err := s.repo.ReadCredential(&transfert.Credential{
		Email: dto.Email,
	})

	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
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
func (s *ClientService) sendMail(credential *entities.Credential, validation *entities.Validation, templateName string) error {
	tpl := template.NewTemplate(templateName)
	if tpl == nil {
		return fmt.Errorf(errors.ErrTemplateNotFound, templateName)
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

	return fmt.Errorf(errors.ErrMailSendFailed)
}

// sendValidationMail Send a signup confirmation email to a client
// This function sends a signup confirmation email to the specified client.
//
// Parameters:
// - client: *entities.Credential The client to send the email to.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *ClientService) sendValidationMail(credential *entities.Credential, token *entities.Validation) error {
	if credential == nil || credential.Email == nil || token == nil {
		// Évitez d'envoyer un e-mail si les données sont manquantes
		return fmt.Errorf("missing data")
	}

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
func (s *ClientService) validateClientAndValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, error) {
	if dtoValidation == nil || dtoCredential == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	credential, err := s.repo.ReadCredential(dtoCredential)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	dtoValidation.ClientID = credential.ClientID
	validation, err := s.repo.ReadValidation(dtoValidation)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrValidationNotFound)
	}

	if validation.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf(errors.ErrValidationExpired)
	}

	if validation.Validated {
		return nil, fmt.Errorf(errors.ErrValidationAlreadyValidated)
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
func (s *ClientService) PasswordValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, error) {
	return s.validateClientAndValidation(dtoValidation, dtoCredential)
}

// SignValidation Validate sign-up
// This function validates a sign-up request.
//
// Parameters:
// - dtoValidation: *transfert.Validation The validation DTO.
// - dtoClient: *transfert.Client The client DTO.
//
// Returns:
// - validation: *entities.Validation The validated validation entity.
// - error: error An error object if an error occurs, nil otherwise.
func (s *ClientService) SignValidation(dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (*entities.Validation, error) {
	return s.validateClientAndValidation(dtoValidation, dtoCredential)
}
