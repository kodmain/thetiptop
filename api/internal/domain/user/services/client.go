package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *UserService) RegisterClient(dtoCredential *transfert.Credential, dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	if dtoCredential == nil || dtoClient == nil {
		return nil, errors.ErrNoDto
	}

	_, err := s.repo.ReadCredential(dtoCredential)
	if err == nil {
		return nil, errors.ErrClientAlreadyExists
	}

	credential, err := s.repo.CreateCredential(dtoCredential)
	if err != nil {
		return nil, errors.FromErr(err, errors.ErrInternalServer)
	}

	dtoClient.CredentialID = &credential.ID

	client, err := s.repo.CreateClient(dtoClient)
	if err != nil {
		return nil, errors.FromErr(err, errors.ErrInternalServer)
	}

	client.Validations = append(client.Validations, &entities.Validation{
		ClientID: &client.ID,
		Type:     entities.MailValidation,
	})

	if err := s.repo.UpdateClient(client); err != nil {
		return nil, errors.FromErr(err, errors.ErrInternalServer)
	}

	if err := s.repo.UpdateCredential(credential); err != nil {
		return nil, errors.FromErr(err, errors.ErrInternalServer)
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
		return nil, errors.ErrClientNotFound
	}

	if !s.security.CanUpdate(client) {
		return nil, errors.ErrUnauthorized
	}

	data.UpdateEntityWithDto(client, dtoClient)

	if err := s.repo.UpdateClient(client); err != nil {
		return nil, errors.FromErr(err, errors.ErrInternalServer)
	}

	return client, nil
}

func (s *UserService) DeleteClient(dtoClient *transfert.Client) errors.ErrorInterface {
	if dtoClient == nil {
		return errors.ErrNoDto
	}

	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return errors.ErrClientNotFound
	}

	if !s.security.CanDelete(client) {
		return errors.ErrUnauthorized
	}

	if err := s.repo.DeleteClient(dtoClient); err != nil {
		return errors.FromErr(err, errors.ErrInternalServer)
	}

	return nil
}

func (s *UserService) GetClient(dtoClient *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	if dtoClient == nil {
		return nil, errors.ErrNoDto
	}

	client, err := s.repo.ReadClient(dtoClient)
	if err != nil {
		return nil, errors.ErrClientNotFound
	}

	if !s.security.CanRead(client) {
		return nil, errors.ErrUnauthorized
	}

	return client, nil
}
