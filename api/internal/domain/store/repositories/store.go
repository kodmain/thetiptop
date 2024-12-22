package repositories

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	errors_domain_store "github.com/kodmain/thetiptop/api/internal/domain/store/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

type StoreRepository struct {
	store *database.Database
}

type StoreRepositoryInterface interface {
	CreateStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface
	ReadStores(obj *transfert.Store, options ...database.Option) ([]*entities.Store, errors.ErrorInterface)
	ReadStore(obj *transfert.Store, options ...database.Option) (*entities.Store, errors.ErrorInterface)
	DeleteStores(obj []*transfert.Store, options ...database.Option) errors.ErrorInterface
	UpdateStores(obj []*entities.Store, options ...database.Option) errors.ErrorInterface

	CreateCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface)
	ReadCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface)
	ReadCaisses(obj *transfert.Caisse, options ...database.Option) ([]*entities.Caisse, errors.ErrorInterface)
	DeleteCaisse(obj *transfert.Caisse, options ...database.Option) errors.ErrorInterface
	UpdateCaisse(obj *entities.Caisse, options ...database.Option) errors.ErrorInterface
}

func NewStoreRepository(repo *database.Database) *StoreRepository {
	repo.Engine.AutoMigrate(entities.Store{}, entities.Caisse{})
	return &StoreRepository{repo}
}

func (r *StoreRepository) CreateStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	stores := make([]*entities.Store, len(objs))
	for i, obj := range objs {
		stores[i] = entities.CreateStore(obj)
	}

	result := r.store.Engine.CreateInBatches(stores, len(stores))
	for _, option := range options {
		option(result)
	}

	if result.Error != nil {
		return errors.ErrInternalServer.Log(result.Error)
	}

	return nil
}

func (r *StoreRepository) ReadStores(obj *transfert.Store, options ...database.Option) ([]*entities.Store, errors.ErrorInterface) {
	var stores []*entities.Store

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.Find(&stores)

	if result.Error != nil {
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return stores, nil
}

func (r *StoreRepository) ReadStore(obj *transfert.Store, options ...database.Option) (*entities.Store, errors.ErrorInterface) {
	var store *entities.Store

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.First(&store)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_store.ErrStoreNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return store, nil
}

func (r *StoreRepository) DeleteStores(obj []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	result := r.store.Engine.Where(obj).Delete(&entities.Store{})

	if result.Error != nil {
		return errors.ErrInternalServer.Log(result.Error)
	}

	return nil
}

func (r *StoreRepository) UpdateStores(obj []*entities.Store, options ...database.Option) errors.ErrorInterface {
	for _, store := range obj {
		result := r.store.Engine.Updates(store)
		for _, option := range options {
			option(result)
		}

		if result.Error != nil {
			return errors.ErrInternalServer.Log(result.Error)
		}
	}

	return nil
}

func (r *StoreRepository) CreateCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	caisse := entities.CreateCaisse(obj)

	result := r.store.Engine.Create(caisse)
	for _, option := range options {
		option(result)
	}

	if result.Error != nil {
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return caisse, nil
}

func (r *StoreRepository) ReadCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	var caisse *entities.Caisse

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.First(&caisse)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_store.ErrCaisseNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return caisse, nil
}

func (r *StoreRepository) ReadCaisses(obj *transfert.Caisse, options ...database.Option) ([]*entities.Caisse, errors.ErrorInterface) {
	var caisses []*entities.Caisse

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.Find(&caisses)

	if result.Error != nil {
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return caisses, nil
}

func (r *StoreRepository) DeleteCaisse(obj *transfert.Caisse, options ...database.Option) errors.ErrorInterface {
	result := r.store.Engine.Where(obj).Delete(&entities.Caisse{})

	if result.Error != nil {
		return errors.ErrInternalServer.Log(result.Error)
	}

	return nil
}

func (r *StoreRepository) UpdateCaisse(obj *entities.Caisse, options ...database.Option) errors.ErrorInterface {
	result := r.store.Engine.Updates(obj)
	for _, option := range options {
		option(result)
	}

	if result.Error != nil {
		return errors.ErrInternalServer.Log(result.Error)
	}

	return nil
}
