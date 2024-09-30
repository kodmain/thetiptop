package repositories

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
)

var user sync.Once

type UserRepository struct {
	store *database.Database
}

type UserRepositoryInterface interface {
	// user
	ReadUser(obj *transfert.User) (*entities.Client, *entities.Employee, errors.ErrorInterface)

	// client
	CreateClient(obj *transfert.Client) (*entities.Client, errors.ErrorInterface)
	ReadClient(obj *transfert.Client) (*entities.Client, errors.ErrorInterface)
	UpdateClient(entity *entities.Client) errors.ErrorInterface
	DeleteClient(obj *transfert.Client) errors.ErrorInterface

	// employee
	CreateEmployee(obj *transfert.Employee) (*entities.Employee, errors.ErrorInterface)
	ReadEmployee(obj *transfert.Employee) (*entities.Employee, errors.ErrorInterface)
	UpdateEmployee(entity *entities.Employee) errors.ErrorInterface
	DeleteEmployee(obj *transfert.Employee) errors.ErrorInterface

	// validation
	CreateValidation(obj *transfert.Validation) (*entities.Validation, errors.ErrorInterface)
	ReadValidation(obj *transfert.Validation) (*entities.Validation, errors.ErrorInterface)
	UpdateValidation(entity *entities.Validation) errors.ErrorInterface
	DeleteValidation(obj *transfert.Validation) errors.ErrorInterface

	// Credential
	CreateCredential(obj *transfert.Credential) (*entities.Credential, errors.ErrorInterface)
	ReadCredential(obj *transfert.Credential) (*entities.Credential, errors.ErrorInterface)
	UpdateCredential(entity *entities.Credential) errors.ErrorInterface
	DeleteCredential(obj *transfert.Credential) errors.ErrorInterface
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

func (r *UserRepository) ReadUser(obj *transfert.User) (*entities.Client, *entities.Employee, errors.ErrorInterface) {
	var client entities.Client
	var employee entities.Employee

	// Chercher dans la table des clients
	resultClient := r.store.Engine.Where(obj.ToClient()).First(&client)
	if resultClient.Error == nil {
		// Si un client est trouvé, le retourner sans chercher d'employé
		return &client, nil, nil
	}

	// Si aucun client trouvé, chercher dans la table des employés
	resultEmployee := r.store.Engine.Where(obj.ToClient()).First(&employee)
	if resultEmployee.Error == nil {
		// Si un employé est trouvé, le retourner
		return nil, &employee, nil
	}

	// Si aucune des deux entités n'est trouvée, retourner une erreur
	if resultClient.Error != nil && resultEmployee.Error != nil {
		return nil, nil, errors_domain_user.ErrUserNotFound
	}

	return nil, nil, errors.ErrInternalServer
}

func (r *UserRepository) CreateCredential(obj *transfert.Credential) (*entities.Credential, errors.ErrorInterface) {
	credential := entities.CreateCredential(obj)

	password, err := hash.Hash(aws.String(*obj.Email+":"+*obj.Password), hash.BCRYPT)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	credential.Password = password

	result := r.store.Engine.Create(credential)
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: credentials.email" {
			return nil, errors_domain_user.ErrCredentialAlreadyExists
		}

		return nil, errors.ErrInternalServer
	}

	return credential, nil
}

func (r *UserRepository) ReadCredential(obj *transfert.Credential) (*entities.Credential, errors.ErrorInterface) {
	credential := entities.CreateCredential(obj)
	result := r.store.Engine.Where(obj).First(credential)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrCredentialNotFound
		}

		return nil, errors.ErrInternalServer
	}

	return credential, nil
}

func (r *UserRepository) UpdateCredential(entity *entities.Credential) errors.ErrorInterface {
	if result := r.store.Engine.Save(entity); result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) DeleteCredential(obj *transfert.Credential) errors.ErrorInterface {
	credential := entities.CreateCredential(obj)
	result := r.store.Engine.Where(obj).Delete(credential)

	if result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) CreateClient(obj *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	client := entities.CreateClient(obj)

	result := r.store.Engine.Create(client)
	if result.Error != nil {
		return nil, errors.ErrInternalServer
	}

	return client, nil
}

func (r *UserRepository) ReadClient(obj *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	client := &entities.Client{}
	result := r.store.Engine.Where(obj).First(client)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrClientNotFound
		}

		return nil, errors.ErrInternalServer
	}

	return client, nil
}

func (r *UserRepository) UpdateClient(entity *entities.Client) errors.ErrorInterface {
	if result := r.store.Engine.Save(entity); result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) DeleteClient(obj *transfert.Client) errors.ErrorInterface {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Delete(client)

	if result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) CreateValidation(obj *transfert.Validation) (*entities.Validation, errors.ErrorInterface) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Create(validation)

	if result.Error != nil {
		return nil, errors.ErrInternalServer
	}

	return validation, nil
}

func (r *UserRepository) ReadValidation(obj *transfert.Validation) (*entities.Validation, errors.ErrorInterface) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Where(obj).First(validation)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrValidationNotFound
		}

		return nil, errors.ErrInternalServer
	}

	return validation, nil
}

func (r *UserRepository) UpdateValidation(entity *entities.Validation) errors.ErrorInterface {
	if result := r.store.Engine.Save(entity); result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) DeleteValidation(obj *transfert.Validation) errors.ErrorInterface {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Where(obj).Delete(validation)

	if result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) CreateEmployee(obj *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	employee := entities.CreateEmployee(obj)

	result := r.store.Engine.Create(employee)
	if result.Error != nil {
		return nil, errors.ErrInternalServer
	}

	return employee, nil
}

func (r *UserRepository) ReadEmployee(obj *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	employee := entities.CreateEmployee(obj)
	result := r.store.Engine.Where(obj).First(employee)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_user.ErrEmployeeNotFound
		}
		return nil, errors.ErrInternalServer
	}

	return employee, nil
}

func (r *UserRepository) UpdateEmployee(entity *entities.Employee) errors.ErrorInterface {
	if result := r.store.Engine.Save(entity); result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) DeleteEmployee(obj *transfert.Employee) errors.ErrorInterface {
	employee := entities.CreateEmployee(obj)
	result := r.store.Engine.Where(obj).Delete(employee)

	if result.Error != nil {
		return errors.ErrInternalServer
	}

	return nil
}
