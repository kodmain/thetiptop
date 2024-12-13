package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *StoreService) ListStores() ([]*entities.Store, errors.ErrorInterface) {
	stores, err := s.repo.ReadStores(&transfert.Store{})
	if err != nil {
		return nil, errors.ErrNoData
	}

	return stores, nil
}
