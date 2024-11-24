package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *UserService) RegisterClient(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	if dtoCredential == nil || dtoClient == nil {
		return nil, errors.ErrNoDto
	}

	_, err := s.repo.ReadCredential(dtoCredential)
	if err == nil {
		return nil, errors_domain_user.ErrClientAlreadyExists
	}

	credential, err := s.repo.CreateCredential(dtoCredential)
	if err != nil {
		return nil, err
	}

	dtoClient.CredentialID = &credential.ID

	client, err := s.repo.CreateClient(dtoClient)
	if err != nil {
		return nil, err
	}

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

func (s *UserService) UpdateClient(dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	if dtoClient == nil {
		return nil, errors.ErrNoDto
	}

	client, err := s.repo.ReadClient(&transfert.Client{
		ID: dtoClient.ID,
	})

	if err != nil {
		return nil, err
	}

	if !s.security.CanUpdate(client) {
		return nil, errors.ErrUnauthorized
	}

	data.UpdateEntityWithDto(client, dtoClient)

	if err := s.repo.UpdateClient(client); err != nil {
		return nil, err
	}

	return client, nil
}

func (s *UserService) DeleteClient(dtoClient *transfert.Client) errors.ErrorInterface {
	if dtoClient == nil {
		return errors.ErrNoDto
	}

	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return err
	}

	if !s.security.CanDelete(client) {
		return errors.ErrUnauthorized
	}

	if err := s.repo.DeleteClient(dtoClient); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetClient(dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	if dtoClient == nil {
		return nil, errors.ErrNoDto
	}

	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return nil, err
	}

	if !s.security.CanRead(client) {
		return nil, errors.ErrUnauthorized
	}

	return client, nil
}
