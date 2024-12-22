package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	user "github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *StoreService) ListStores() ([]*entities.Store, errors.ErrorInterface) {
	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	stores, err := s.repo.ReadStores(&transfert.Store{})
	if err != nil {
		return nil, errors.ErrNoData
	}

	return stores, nil
}

func (s *StoreService) GetStoreByID(dto *transfert.Store) (*entities.Store, errors.ErrorInterface) {
	if dto == nil {
		return nil, errors.ErrNoDto
	}

	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	store, err := s.repo.ReadStore(dto)
	if err != nil {
		return nil, err
	}

	return store, nil
}
