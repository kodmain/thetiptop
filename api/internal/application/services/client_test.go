package services_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	email              = "test@example.com"
	emailSyntaxFail    = "testexample.com"
	password           = "validP@ssw0rd"
	passwordFail       = "WrongP@ssw0rd"
	passwordSyntaxFail = "secret"
	trueValue          = aws.Bool(true)
	falseValue         = aws.Bool(false)

	ExpiredAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDgyMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJyZWZyZXNoIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM01UTXhNRGt4TXpFc0ltbGtJam9pTjJNM09UUXdNR1l0TURBMllTMDBOelZsTFRrM1lqWXROV1JpWkdVek56QTNOakF4SWl3aWIyWm1Jam8zTWpBd0xDSjBlWEJsSWpveExDSjBlaUk2SWt4dlkyRnNJbjAuNUxhZTU2SE5jUTFPSGNQX0ZoVGZjT090SHBhWlZnUkZ5NnZ6ekJ1Z043WSIsInR5cGUiOjAsInR6IjoiTG9jYWwifQ.BxW2wfHiiCr0aTsuWwRVmh0Wd-BX20AoUDTGg_rIDoM"
	ExpiredRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y"
)

type DomainClientService struct {
	mock.Mock
	mu sync.Mutex
}

func (dcs *DomainClientService) GetClient(client *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainClientService) PasswordRecover(obj *transfert.Credential) error {
	args := dcs.Called(obj)
	return args.Error(0)
}

func (dcs *DomainClientService) UserRegister(credential *transfert.Credential, client *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(credential, client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainClientService) UserAuth(obj *transfert.Credential) (*entities.Client, error) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs *DomainClientService) MailValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, error) {
	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs *DomainClientService) ValidationRecover(validation *transfert.Validation, credential *transfert.Credential) error {
	args := dcs.Called(validation, credential)
	return args.Error(0)
}

func (dcs *DomainClientService) PasswordValidation(validation *transfert.Validation, credential *transfert.Credential) (*entities.Validation, error) {
	dcs.mu.Lock()
	defer dcs.mu.Unlock()

	args := dcs.Called(validation, credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs *DomainClientService) UpdateClient(client *transfert.Client) error {
	args := dcs.Called(client)
	return args.Error(0)
}

func (dcs *DomainClientService) DeleteClient(client *transfert.Client) error {
	args := dcs.Called(client)
	return args.Error(0)
}

func (dcs *DomainClientService) PasswordUpdate(credential *transfert.Credential) error {
	args := dcs.Called(credential)
	return args.Error(0)
}

func TestMailValidation(t *testing.T) {
	luhn := token.Generate(6)

	// Cas de validation réussie
	t.Run("successful validation", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Génération d'un ID aléatoire pour la validation
		id, err := uuid.NewRandom()
		assert.NoError(t, err)

		// Simulation de la réponse du mock pour une validation réussie
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(&entities.Validation{
			ID: id.String(),
		}, nil)

		// Appel de la fonction à tester
		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})

		// Vérifications
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
		assert.IsType(t, &entities.Validation{}, response)

		// Vérification des attentes du mock
		mockClient.AssertExpectations(t)
	})

	// Cas d'un token avec une syntaxe invalide
	t.Run("invalid token syntax", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Appel avec un token invalide
		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: aws.String("invalidToken"),
		}, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		})

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "invalid digit", response) // Correction du message d'erreur attendu
	})

	// Cas d'un token avec une syntaxe invalide
	t.Run("invalid email syntax", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Appel avec un token invalide
		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		})

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "mail: missing '@' or angle-addr", response) // Correction du message d'erreur attendu
	})

	// Cas où la validation n'a pas été trouvée
	t.Run("validation not found", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Configuration du mock pour renvoyer l'erreur "validation not found"
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrValidationNotFound, response)
		mockClient.AssertExpectations(t) // Vérification des attentes du mock
	})

	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationExpired))

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, errors.ErrValidationExpired, response.(string))
	})

	// Cas où la validation a déjà été effectuée
	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Configuration du mock pour renvoyer l'erreur "validation already validated"
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationAlreadyValidated))

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})

		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.Equal(t, errors.ErrValidationAlreadyValidated, response)
		mockClient.AssertExpectations(t)
	})
}

func TestValidationRecover(t *testing.T) {
	validationType := "password_recovery"

	t.Run("invalid email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})
		// Ajustement du message d'erreur pour correspondre à celui renvoyé par le validateur
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "mail: missing '@' or angle-addr", response)
	})

	t.Run("invalid email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: nil,
		})
		// Ajustement du message d'erreur pour correspondre à celui renvoyé par le validateur
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value type is required", response)
	})

	t.Run("successful recovery", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Mock pour un cas de récupération réussi
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).Return(nil)

		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})
		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("client not found", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Mock pour simuler un client non trouvé
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).Return(fmt.Errorf(errors.ErrClientNotFound))

		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})

		// Correction des attentes pour vérifier le bon statut et message
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrClientNotFound, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("other error", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Mock pour simuler un client non trouvé
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).Return(fmt.Errorf("random error"))

		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})

		// Correction des attentes pour vérifier le bon statut et message
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "random error", response)
		mockClient.AssertExpectations(t)
	})
}

func TestUpdateClient(t *testing.T) {

	t.Run("invalid client data", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Test pour vérifier la validation de l'ID (UUID)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         nil, // ID manquant
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value id is required", response) // Ajustement du message attendu
	})

	t.Run("invalid newsletter data", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Test pour vérifier la validation de la valeur newsletter
		// Aucune interaction avec UpdateClient ici, car la validation échoue avant d'appeler cette méthode
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Valeur incorrecte pour Newsletter
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value newsletter is required", response) // Message attendu pour erreur booléenne
	})

	t.Run("successful client update", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Mock pour simuler un cas de mise à jour réussie
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(nil)

		// Fournir un UUID valide ici
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("client update error", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Simuler une erreur lors de la mise à jour du client
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(fmt.Errorf("update error"))

		// Fournir un UUID valide mais avec un champ Newsletter manquant (incomplet)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Manque une valeur pour `Newsletter`
		})
		// La validation va échouer avant que `UpdateClient` ne soit appelé, donc nous testons cela
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value newsletter is required", response) // Message d'erreur attendu
		mockClient.AssertNotCalled(t, "UpdateClient")             // Vérifier que le mock n'a pas été appelé
	})

	t.Run("client update error", func(t *testing.T) {
		mockClient := new(DomainClientService)
		// Simuler une erreur lors de la mise à jour du client
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(fmt.Errorf("update error"))

		// Fournir un UUID valide mais avec un champ Newsletter manquant (incomplet)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: aws.Bool(true),                                     // Manque une valeur pour `Newsletter`
		})
		// La validation va échouer avant que `UpdateClient` ne soit appelé, donc nous testons cela
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, "update error", response)     // Message d'erreur attendu
		mockClient.AssertNotCalled(t, "UpdateClient") // Vérifier que le mock n'a pas été appelé
	})

}

func TestGetClient(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		mockService := new(DomainClientService)
		// Simuler une entrée incorrecte avec un ID manquant
		dtoClient := &transfert.Client{
			ID: nil, // Manque l'ID
		}

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value id is required", response)
	})

	t.Run("client not found", func(t *testing.T) {
		mockService := new(DomainClientService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du service qui ne trouve pas le client
		mockService.On("GetClient", dtoClient).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrClientNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockService := new(DomainClientService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler une erreur interne inattendue
		mockService.On("GetClient", dtoClient).Return(nil, fmt.Errorf("internal error"))

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, "internal error", response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful client retrieval", func(t *testing.T) {
		mockService := new(DomainClientService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du service avec un client valide
		expectedClient := &entities.Client{
			ID:  "42debee6-2063-4566-baf1-37a7bdd139ff",
			CGU: aws.Bool(true),
		}

		mockService.On("GetClient", dtoClient).Return(expectedClient, nil)

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedClient, response)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteClient(t *testing.T) {
	t.Run("should return 400 if validation fails", func(t *testing.T) {
		// Créer un mock pour le service client
		mockService := new(DomainClientService)

		// Cas où l'ID du client est manquant
		dtoClient := &transfert.Client{ID: nil}

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "value id is required", response)
	})

	t.Run("should return 404 if client not found", func(t *testing.T) {
		// Créer un mock pour le service client
		mockService := new(DomainClientService)

		// Cas où le client n'est pas trouvé
		clientID := "123e4567-e89b-12d3-a456-426614174000"
		dtoClient := &transfert.Client{ID: &clientID}

		// Configurer le mock pour renvoyer une erreur client non trouvé
		mockService.On("DeleteClient", dtoClient).Return(fmt.Errorf(errors.ErrClientNotFound))

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrClientNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 if internal server error occurs", func(t *testing.T) {

		// Créer un mock pour le service client
		mockService := new(DomainClientService)
		// Cas où une erreur interne survient lors de la suppression
		clientID := "123e4567-e89b-12d3-a456-426614174000"
		dtoClient := &transfert.Client{ID: &clientID}

		// Configurer le mock pour renvoyer une erreur interne
		mockService.On("DeleteClient", dtoClient).Return(fmt.Errorf("internal error"))

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, "internal error", response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 204 if client is deleted successfully", func(t *testing.T) {
		// Créer un mock pour le service client
		mockService := new(DomainClientService)
		// Cas de suppression réussie
		clientID := "123e4567-e89b-12d3-a456-426614174000"
		dtoClient := &transfert.Client{ID: &clientID}

		// Configurer le mock pour ne pas renvoyer d'erreur
		mockService.On("DeleteClient", dtoClient).Return(nil)

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
		mockService.AssertExpectations(t)
	})
}
