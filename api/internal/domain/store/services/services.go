package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type StoreService struct {
	security security.PermissionInterface
	repo     repositories.StoreRepositoryInterface
}

func Store(security security.PermissionInterface, repo repositories.StoreRepositoryInterface) *StoreService {
	return &StoreService{security, repo}
}

type StoreServiceInterface interface {
	ListStores() ([]*entities.Store, errors.ErrorInterface)
	GetStoreByID(*transfert.Store) (*entities.Store, errors.ErrorInterface)

	GetCaisse(*transfert.Caisse) (*entities.Caisse, errors.ErrorInterface)
	CreateCaisse(*transfert.Caisse) (*entities.Caisse, errors.ErrorInterface)
	DeleteCaisse(*transfert.Caisse) errors.ErrorInterface
	UpdateCaisse(*transfert.Caisse) (*entities.Caisse, errors.ErrorInterface)
}
