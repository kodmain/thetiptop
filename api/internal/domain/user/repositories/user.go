package repositories

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
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
	ReadUser(obj *transfert.User) (*entities.Client, *entities.Employee, error)

	// client
	CreateClient(obj *transfert.Client) (*entities.Client, error)
	ReadClient(obj *transfert.Client) (*entities.Client, error)
	UpdateClient(entity *entities.Client) error
	DeleteClient(obj *transfert.Client) error

	// employee
	CreateEmployee(obj *transfert.Employee) (*entities.Employee, error)
	ReadEmployee(obj *transfert.Employee) (*entities.Employee, error)
	UpdateEmployee(entity *entities.Employee) error
	DeleteEmployee(obj *transfert.Employee) error

	// validation
	CreateValidation(obj *transfert.Validation) (*entities.Validation, error)
	ReadValidation(obj *transfert.Validation) (*entities.Validation, error)
	UpdateValidation(entity *entities.Validation) error
	DeleteValidation(obj *transfert.Validation) error

	// Credential
	CreateCredential(obj *transfert.Credential) (*entities.Credential, error)
	ReadCredential(obj *transfert.Credential) (*entities.Credential, error)
	UpdateCredential(entity *entities.Credential) error
	DeleteCredential(obj *transfert.Credential) error
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

func (r *UserRepository) ReadUser(obj *transfert.User) (*entities.Client, *entities.Employee, error) {
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
		return nil, nil, errors.ErrUserNotFound
	}

	return nil, nil, errors.ErrUnknown
}

func (r *UserRepository) CreateCredential(obj *transfert.Credential) (*entities.Credential, error) {
	credential := entities.CreateCredential(obj)

	password, err := hash.Hash(aws.String(*obj.Email+":"+*obj.Password), hash.BCRYPT)
	if err != nil {
		return nil, err
	}

	credential.Password = password

	result := r.store.Engine.Create(credential)
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: credentials.email" {
			return nil, errors.ErrCredentialAlreadyExists
		}

		return nil, result.Error
	}

	return credential, nil
}

func (r *UserRepository) ReadCredential(obj *transfert.Credential) (*entities.Credential, error) {
	credential := entities.CreateCredential(obj)
	result := r.store.Engine.Where(obj).First(credential)

	if result.Error != nil {
		return nil, result.Error
	}

	return credential, nil
}

func (r *UserRepository) UpdateCredential(entity *entities.Credential) error {
	return r.store.Engine.Save(entity).Error
}

func (r *UserRepository) DeleteCredential(obj *transfert.Credential) error {
	credential := entities.CreateCredential(obj)
	result := r.store.Engine.Where(obj).Delete(credential)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) CreateClient(obj *transfert.Client) (*entities.Client, error) {
	client := entities.CreateClient(obj)

	result := r.store.Engine.Create(client)
	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}

func (r *UserRepository) ReadClient(obj *transfert.Client) (*entities.Client, error) {
	client := &entities.Client{}
	result := r.store.Engine.Where(obj).First(client)

	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}

func (r *UserRepository) UpdateClient(entity *entities.Client) error {
	return r.store.Engine.Save(entity).Error
}

func (r *UserRepository) DeleteClient(obj *transfert.Client) error {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Delete(client)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) CreateValidation(obj *transfert.Validation) (*entities.Validation, error) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Create(validation)

	if result.Error != nil {
		return nil, result.Error
	}

	return validation, nil
}

func (r *UserRepository) ReadValidation(obj *transfert.Validation) (*entities.Validation, error) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Where(obj).First(validation)

	if result.Error != nil {
		return nil, result.Error
	}

	return validation, nil
}

func (r *UserRepository) UpdateValidation(entity *entities.Validation) error {
	return r.store.Engine.Save(entity).Error
}

func (r *UserRepository) DeleteValidation(obj *transfert.Validation) error {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Where(obj).Delete(validation)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) CreateEmployee(obj *transfert.Employee) (*entities.Employee, error) {
	employee := entities.CreateEmployee(obj)

	result := r.store.Engine.Create(employee)
	if result.Error != nil {
		return nil, result.Error
	}

	return employee, nil
}

func (r *UserRepository) ReadEmployee(obj *transfert.Employee) (*entities.Employee, error) {
	employee := entities.CreateEmployee(obj)
	result := r.store.Engine.Where(obj).First(employee)

	if result.Error != nil {
		return nil, result.Error
	}

	return employee, nil
}

func (r *UserRepository) UpdateEmployee(entity *entities.Employee) error {
	return r.store.Engine.Save(entity).Error
}

func (r *UserRepository) DeleteEmployee(obj *transfert.Employee) error {
	employee := entities.CreateEmployee(obj)
	result := r.store.Engine.Where(obj).Delete(employee)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
