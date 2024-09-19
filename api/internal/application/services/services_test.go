package services_test

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
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

func (dcs *DomainUserService) GetClient(client *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainUserService) PasswordRecover(obj *transfert.Credential) error {
	args := dcs.Called(obj)
	return args.Error(0)
}

func (dcs *DomainUserService) UserAuth(obj *transfert.Credential) (*string, error) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Retourne nil pour *string et l'erreur s'il y en a une
	}

	return args.Get(0).(*string), args.Error(1)
}

func (dcs *DomainUserService) MailValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, error) {
	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs *DomainUserService) ValidationRecover(validation *transfert.Validation, credential *transfert.Credential) error {
	args := dcs.Called(validation, credential)
	return args.Error(0)
}

func (dcs *DomainUserService) PasswordValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, error) {
	dcs.mu.Lock()
	defer dcs.mu.Unlock()

	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs *DomainUserService) UpdateClient(client *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainUserService) DeleteClient(client *transfert.Client) error {
	args := dcs.Called(client)
	return args.Error(0)
}

func (dcs *DomainUserService) PasswordUpdate(credential *transfert.Credential) error {
	args := dcs.Called(credential)
	return args.Error(0)
}

func (dcs *DomainUserService) RegisterClient(credential *transfert.Credential, client *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(credential, client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainUserService) RegisterEmployee(credential *transfert.Credential, employee *transfert.Employee) (*entities.Employee, error) {
	args := dcs.Called(credential, employee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Employee), args.Error(1)
}

func (dcs *DomainUserService) GetEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, error) {
	args := dcs.Called(dtoEmployee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Employee), args.Error(1)
}

func (dcs *DomainUserService) DeleteEmployee(dtoEmployee *transfert.Employee) error {
	args := dcs.Called(dtoEmployee)
	return args.Error(0)
}

func (dcs *DomainUserService) UpdateEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, error) {
	args := dcs.Called(dtoEmployee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Employee), args.Error(1)
}
