package services

import (
	"fmt"
	"time"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
)

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
func (s *ClientService) validateClientAndValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	dtoValidation.ClientID = &client.ID
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
func (s *ClientService) PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	return s.validateClientAndValidation(dtoValidation, dtoClient)
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
func (s *ClientService) SignValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	return s.validateClientAndValidation(dtoValidation, dtoClient)
}
