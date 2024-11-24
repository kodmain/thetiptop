package repositories

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"gorm.io/gorm"
)

var user sync.Once

type UserRepository struct {
	store *database.Database
}

type UserRepositoryInterface interface {
	// user
	ReadUser(obj *transfert.User, options ...database.Option) (*entities.Client, *entities.Employee, errors.ErrorInterface)

	// client
	CreateClient(obj *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface)
	ReadClient(obj *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface)
	UpdateClient(entity *entities.Client, options ...database.Option) errors.ErrorInterface
	DeleteClient(obj *transfert.Client, options ...database.Option) errors.ErrorInterface

	// employee
	CreateEmployee(obj *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface)
	ReadEmployee(obj *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface)
	UpdateEmployee(entity *entities.Employee, options ...database.Option) errors.ErrorInterface
	DeleteEmployee(obj *transfert.Employee, options ...database.Option) errors.ErrorInterface

	// validation
	CreateValidation(obj *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface)
	ReadValidation(obj *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface)
	UpdateValidation(entity *entities.Validation, options ...database.Option) errors.ErrorInterface
	DeleteValidation(obj *transfert.Validation, options ...database.Option) errors.ErrorInterface

	// Credential
	CreateCredential(obj *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface)
	ReadCredential(obj *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface)
	UpdateCredential(entity *entities.Credential, options ...database.Option) errors.ErrorInterface
	DeleteCredential(obj *transfert.Credential, options ...database.Option) errors.ErrorInterface
}

func NewUserRepository(store *database.Database) *UserRepository {
	user.Do(func() {
		store.Engine.AutoMigrate(entities.Client{})
		store.Engine.AutoMigrate(entities.Employee{})
		store.Engine.AutoMigrate(entities.Validation{})
		store.Engine.AutoMigrate(entities.Credential{})
	})

	return &UserRepository{store}
}

func (r *UserRepository) applyOptions(query *gorm.DB, options ...database.Option) {
	for _, option := range options {
		option(query)
	}
}

func (r *UserRepository) ReadUser(obj *transfert.User, options ...database.Option) (*entities.Client, *entities.Employee, errors.ErrorInterface) {
	var client entities.Client
	var employee entities.Employee

	queryClient := r.store.Engine.Where(obj.ToClient())
	r.applyOptions(queryClient, options...)
	resultClient := queryClient.First(&client)
	if resultClient.Error == nil {
		return &client, nil, nil
	}

	queryEmployee := r.store.Engine.Where(obj.ToEmployee())
	r.applyOptions(queryEmployee, options...)
	resultEmployee := queryEmployee.First(&employee)
	if resultEmployee.Error == nil {
		return nil, &employee, nil
	}

	if resultClient.Error != nil && resultEmployee.Error != nil {
		return nil, nil, errors_domain_user.ErrUserNotFound
	}

	return nil, nil, errors.ErrInternalServer.Log(resultClient.Error)
}

func (r *UserRepository) CreateCredential(obj *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface) {
	credential := entities.CreateCredential(obj)

	password, err := hash.Hash(aws.String(*obj.Email+":"+*obj.Password), hash.BCRYPT)
	if err != nil {
		return nil, errors.ErrInternalServer.Log(err)
	}

	credential.Password = password

	query := r.store.Engine.Create(credential)
	r.applyOptions(query, options...)
	if query.Error != nil {
		if query.Error.Error() == "UNIQUE constraint failed: credentials.email" {
			return nil, errors_domain_user.ErrCredentialAlreadyExists
		}
		return nil, errors.ErrInternalServer.Log(query.Error)
	}

	return credential, nil
}

func (r *UserRepository) ReadCredential(obj *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface) {
	credential := entities.CreateCredential(obj)
	query := r.store.Engine.Where(obj)
	r.applyOptions(query, options...)
	result := query.First(credential)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrCredentialNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return credential, nil
}

func (r *UserRepository) UpdateCredential(entity *entities.Credential, options ...database.Option) errors.ErrorInterface {
	query := r.store.Engine.Save(entity)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) DeleteCredential(obj *transfert.Credential, options ...database.Option) errors.ErrorInterface {
	credential := entities.CreateCredential(obj)
	query := r.store.Engine.Where(obj).Delete(credential)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) CreateClient(obj *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface) {
	client := entities.CreateClient(obj)

	query := r.store.Engine.Create(client)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return nil, errors.ErrInternalServer.Log(query.Error)
	}

	return client, nil
}

func (r *UserRepository) ReadClient(obj *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface) {
	client := &entities.Client{}
	query := r.store.Engine.Where(obj)
	r.applyOptions(query, options...)
	result := query.First(client)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrClientNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return client, nil
}

func (r *UserRepository) UpdateClient(entity *entities.Client, options ...database.Option) errors.ErrorInterface {
	query := r.store.Engine.Save(entity)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) DeleteClient(obj *transfert.Client, options ...database.Option) errors.ErrorInterface {
	client := entities.CreateClient(obj)
	query := r.store.Engine.Where(obj).Delete(client)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) CreateValidation(obj *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface) {
	validation := entities.CreateValidation(obj)
	query := r.store.Engine.Create(validation)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return nil, errors.ErrInternalServer.Log(query.Error)
	}

	return validation, nil
}

func (r *UserRepository) ReadValidation(obj *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface) {
	validation := entities.CreateValidation(obj)
	query := r.store.Engine.Where(obj)
	r.applyOptions(query, options...)
	result := query.First(validation)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrValidationNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return validation, nil
}

func (r *UserRepository) UpdateValidation(entity *entities.Validation, options ...database.Option) errors.ErrorInterface {
	query := r.store.Engine.Save(entity)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) DeleteValidation(obj *transfert.Validation, options ...database.Option) errors.ErrorInterface {
	validation := entities.CreateValidation(obj)
	query := r.store.Engine.Where(obj).Delete(validation)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) CreateEmployee(obj *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface) {
	employee := entities.CreateEmployee(obj)

	query := r.store.Engine.Create(employee)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return nil, errors.ErrInternalServer.Log(query.Error)
	}

	return employee, nil
}

func (r *UserRepository) ReadEmployee(obj *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface) {
	employee := entities.CreateEmployee(obj)
	query := r.store.Engine.Where(obj)
	r.applyOptions(query, options...)
	result := query.First(employee)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrEmployeeNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return employee, nil
}

func (r *UserRepository) UpdateEmployee(entity *entities.Employee, options ...database.Option) errors.ErrorInterface {
	query := r.store.Engine.Save(entity)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

func (r *UserRepository) DeleteEmployee(obj *transfert.Employee, options ...database.Option) errors.ErrorInterface {
	employee := entities.CreateEmployee(obj)
	query := r.store.Engine.Where(obj).Delete(employee)
	r.applyOptions(query, options...)
	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}
