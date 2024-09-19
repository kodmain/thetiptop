package services

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func (s *UserService) RegisterClient(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, error) {
	if dtoCredential == nil || dtoClient == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

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

	client.CredentialID = &credential.ID

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

func (s *UserService) UpdateClient(dtoClient *transfert.Client) (*entities.Client, error) {
	client, err := s.repo.ReadClient(&transfert.Client{
		ID: dtoClient.ID,
	})

	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	if err := data.UpdateEntityWithDto(client, dtoClient); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateClient(client); err != nil {
		return nil, err
	}

	return client, nil
}

// DeleteClient Delete a client from the repository
// This function deletes a client entity from the repository, using the provided client DTO.
//
// Parameters:
// - dtoClient: *transfert.Client The client DTO containing the ID of the client to delete.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *UserService) DeleteClient(dtoClient *transfert.Client) error {
	// Vérifier si le DTO est valide
	if dtoClient == nil || dtoClient.ID == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	// Supprimer le client du repository
	if err := s.repo.DeleteClient(dtoClient); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetClient(dtoClient *transfert.Client) (*entities.Client, error) {
	if dtoClient == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	return client, nil
}
