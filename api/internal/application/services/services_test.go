package services_test

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/stretchr/testify/mock"
)

var (
	email              = "test@example.com"
	emailSyntaxFail    = "testexample.com"
	password           = "validP@ssw0rd"
	passwordSyntaxFail = "secret"
	trueValue          = aws.Bool(true)
	falseValue         = aws.Bool(false)

	ExpiredAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDgyMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJyZWZyZXNoIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM01UTXhNRGt4TXpFc0ltbGtJam9pTjJNM09UUXdNR1l0TURBMllTMDBOelZsTFRrM1lqWXROV1JpWkdVek56QTNOakF4SWl3aWIyWm1Jam8zTWpBd0xDSjBlWEJsSWpveExDSjBlaUk2SWt4dlkyRnNJbjAuNUxhZTU2SE5jUTFPSGNQX0ZoVGZjT090SHBhWlZnUkZ5NnZ6ekJ1Z043WSIsInR5cGUiOjAsInR6IjoiTG9jYWwifQ.BxW2wfHiiCr0aTsuWwRVmh0Wd-BX20AoUDTGg_rIDoM"
	ExpiredRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y"
)

type DomainUserService struct {
	mock.Mock
	mu sync.Mutex
}

func (dcs *DomainUserService) GetClient(client *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (dcs *DomainUserService) PasswordRecover(obj *transfert.Credential) errors.ErrorInterface {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (dcs *DomainUserService) UserAuth(obj *transfert.Credential) (*string, string, errors.ErrorInterface) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, "", args.Get(2).(errors.ErrorInterface) // Retourne nil pour *string et l'erreur s'il y en a une
	}

	return args.Get(0).(*string), args.Get(1).(string), nil
}

func (dcs *DomainUserService) MailValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, errors.ErrorInterface) {
	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (dcs *DomainUserService) ValidationRecover(validation *transfert.Validation, credential *transfert.Credential) errors.ErrorInterface {
	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (dcs *DomainUserService) PasswordValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, errors.ErrorInterface) {
	dcs.mu.Lock()
	defer dcs.mu.Unlock()

	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (dcs *DomainUserService) UpdateClient(client *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}

	return args.Get(0).(*entities.Client), nil
}

func (dcs *DomainUserService) DeleteClient(client *transfert.Client) errors.ErrorInterface {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (dcs *DomainUserService) PasswordUpdate(credential *transfert.Credential) errors.ErrorInterface {
	args := dcs.Called(credential)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (dcs *DomainUserService) RegisterClient(credential *transfert.Credential, client *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	args := dcs.Called(credential, client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (dcs *DomainUserService) RegisterEmployee(credential *transfert.Credential, employee *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	args := dcs.Called(credential, employee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (dcs *DomainUserService) GetEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	args := dcs.Called(dtoEmployee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (dcs *DomainUserService) DeleteEmployee(dtoEmployee *transfert.Employee) errors.ErrorInterface {
	args := dcs.Called(dtoEmployee)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (dcs *DomainUserService) UpdateEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	args := dcs.Called(dtoEmployee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}
