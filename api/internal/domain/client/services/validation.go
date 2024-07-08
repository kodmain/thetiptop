package services

import (
	"fmt"
	"time"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
)

func (s *ClientService) PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	client, err := s.repo.ReadClient(&transfert.Client{
		Email: dtoClient.Email,
	})

	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	dtoValidation.ClientID = &client.ID
	validation, err := s.repo.ReadValidation(dtoValidation)

	if err != nil {
		return nil, fmt.Errorf(errors.ErrValidationNotFound)
	}

	if validation.Validated {
		return nil, fmt.Errorf(errors.ErrValidationAlreadyValidated)
	}

	if validation.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf(errors.ErrValidationExpired)
	}

	validation.Validated = true

	if s.repo.UpdateValidation(validation) != nil {
		return nil, err
	}

	return validation, nil
}

func (s *ClientService) SignValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	dtoValidation.ClientID = &client.ID
	validation, err := s.repo.ReadValidation(dtoValidation)

	if err != nil {
		return nil, fmt.Errorf(errors.ErrValidationNotFound)
	}

	if validation.Validated {
		return nil, fmt.Errorf(errors.ErrValidationAlreadyValidated)
	}

	if validation.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf(errors.ErrValidationExpired)
	}

	validation.Validated = true

	if s.repo.UpdateValidation(validation) != nil {
		return nil, err
	}

	return validation, nil
}
