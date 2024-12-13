package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	user "github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *StoreService) GetCaisse(dto *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	if dto == nil {
		return nil, errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	caisse, err := s.repo.ReadCaisse(dto)
	if err != nil {
		return nil, err
	}

	return caisse, nil
}

func (s *StoreService) GetCaissesByStore(dto *transfert.Caisse) ([]*entities.Caisse, errors.ErrorInterface) {
	if dto == nil {
		return nil, errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	caisses, err := s.repo.ReadCaisses(dto)
	if err != nil {
		return nil, err
	}

	return caisses, nil
}

func (s *StoreService) CreateCaisse(dto *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	if dto == nil {
		return nil, errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	_, err := s.repo.ReadStore(&transfert.Store{ID: dto.StoreID})
	if err != nil {
		return nil, err
	}

	caisse, err := s.repo.CreateCaisse(dto)
	if err != nil {
		return nil, err
	}

	return caisse, nil
}

func (s *StoreService) UpdateCaisse(dto *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	if dto == nil {
		return nil, errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	caisse, err := s.repo.ReadCaisse(&transfert.Caisse{ID: dto.ID})
	if err != nil {
		return nil, err
	}

	data.UpdateEntityWithDto(caisse, dto)

	if err := s.repo.UpdateCaisse(caisse); err != nil {
		return nil, err
	}

	return caisse, nil
}

func (s *StoreService) DeleteCaisse(dto *transfert.Caisse) errors.ErrorInterface {
	if dto == nil {
		return errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return errors.ErrUnauthorized
	}

	if err := s.repo.DeleteCaisse(dto); err != nil {
		return err
	}

	return nil
}
