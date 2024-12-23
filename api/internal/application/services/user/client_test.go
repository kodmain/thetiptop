package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	services "github.com/kodmain/thetiptop/api/internal/application/services/user"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	entitiesGame "github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterClient(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

	t.Run("invalid password", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mock de la méthode RegisterClient
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(&entities.Client{}, nil)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &passwordSyntaxFail,
		}, &transfert.Client{
			Newsletter: trueValue,
			CGU:        trueValue,
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("missing newsletter", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Pas besoin de mocker RegisterClient car l'erreur survient avant l'appel

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		}, &transfert.Client{
			Newsletter: nil, // Newsletter manquant
			CGU:        trueValue,
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("invalid cgu", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Pas besoin de mocker RegisterClient car l'erreur survient avant l'appel

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		}, &transfert.Client{
			Newsletter: trueValue,
			CGU:        falseValue, // CGU doit être à true
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("valid password and fields", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mock pour simuler un cas de succès
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(&entities.Client{}, nil)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		}, &transfert.Client{
			Newsletter: trueValue,
			CGU:        trueValue,
		})

		assert.Equal(t, fiber.StatusCreated, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("client already exists", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler le cas où le client existe déjà
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(nil, errors_domain_user.ErrCredentialAlreadyExists)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		}, &transfert.Client{
			Newsletter: trueValue,
			CGU:        trueValue,
		})
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("server error during registration", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une erreur serveur lors de la tentative d'enregistrement
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(nil, errors.ErrInternalServer)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		}, &transfert.Client{
			Newsletter: trueValue,
			CGU:        trueValue,
		})
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
	})
}

func TestMailValidation(t *testing.T) {
	luhn := token.Generate(6)

	// Cas de validation réussie
	t.Run("successful validation", func(t *testing.T) {
		mockClient := new(DomainUserService)
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
		mockClient := new(DomainUserService)
		// Appel avec un token invalide
		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: aws.String("invalidToken"),
		}, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		})

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "validator.is_not_number", response) // Correction du message d'erreur attendu
	})

	// Cas d'un token avec une syntaxe invalide
	t.Run("invalid email syntax", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Appel avec un token invalide
		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		})

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Equal(t, "validator.is_not_email", response) // Correction du message d'erreur attendu
	})

	// Cas où la validation n'a pas été trouvée
	t.Run("validation not found", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Configuration du mock pour renvoyer l'erreur "validation not found"
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationNotFound)

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationNotFound, response)
		mockClient.AssertExpectations(t) // Vérification des attentes du mock
	})

	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainUserService)
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationExpired)

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, errors_domain_user.ErrValidationExpired, response)
	})

	// Cas où la validation a déjà été effectuée
	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Configuration du mock pour renvoyer l'erreur "validation already validated"
		mockClient.On("MailValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationAlreadyValidated)

		statusCode, response := services.MailValidation(mockClient, &transfert.Validation{
			Token: luhn.PointerString(),
		}, &transfert.Credential{
			Email: aws.String(email),
		})

		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationAlreadyValidated, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("should return 200 and client data on success", func(t *testing.T) {
		mockService := new(DomainUserService)

		// Simuler des données de client exportées
		expectedClients := &entities.ClientData{
			Credential:  &entities.Credential{ID: "credential-id", Email: aws.String("email@example.com")},
			Client:      &entities.Client{ID: "client-id", CGU: aws.Bool(true)},
			Validations: []*entities.Validation{{ID: "validation-id", Validated: false}},
			Tickets:     []*entitiesGame.Ticket{{ID: "ticket-id"}},
		}

		// Configurer le mock pour retourner des données valides
		mockService.On("ExportClient").Return(expectedClients, nil)

		// Appeler la méthode à tester
		statusCode, response := services.ExportClient(mockService)

		// Vérifications
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedClients, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on internal error", func(t *testing.T) {
		mockService := new(DomainUserService)

		// Simuler une erreur interne
		mockService.On("ExportClient").Return(nil, errors.ErrInternalServer)

		// Appeler la méthode à tester
		statusCode, response := services.ExportClient(mockService)

		// Vérifications
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on unauthorized", func(t *testing.T) {
		mockService := new(DomainUserService)

		// Simuler une erreur d'autorisation
		mockService.On("ExportClient").Return(nil, errors.ErrUnauthorized)

		// Appeler la méthode à tester
		statusCode, response := services.ExportClient(mockService)

		// Vérifications
		assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		assert.Equal(t, errors.ErrUnauthorized, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on client not found", func(t *testing.T) {
		mockService := new(DomainUserService)

		// Simuler une erreur de client introuvable
		mockService.On("ExportClient").Return(nil, errors.ErrNotFound)

		// Appeler la méthode à tester
		statusCode, response := services.ExportClient(mockService)

		// Vérifications
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrNotFound, response)
		mockService.AssertExpectations(t)
	})
}

func TestValidationRecover(t *testing.T) {
	validationType := "password_recovery"

	t.Run("invalid email", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(emailSyntaxFail),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})
		// Ajustement du message d'erreur pour correspondre à celui renvoyé par le validateur
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("invalid email", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: nil,
		})
		// Ajustement du message d'erreur pour correspondre à celui renvoyé par le validateur
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("successful recovery", func(t *testing.T) {
		mockClient := new(DomainUserService)
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
		mockClient := new(DomainUserService)
		// Mock pour simuler un client non trouvé
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).Return(errors_domain_user.ErrUserNotFound)

		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})

		// Correction des attentes pour vérifier le bon statut et message
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrUserNotFound, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("other error", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mock pour simuler un client non trouvé
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).Return(errors.ErrInternalServer)

		statusCode, response := services.ValidationRecover(mockClient, &transfert.Credential{
			Email: aws.String(email),
		}, &transfert.Validation{
			Type: aws.String(validationType),
		})

		// Correction des attentes pour vérifier le bon statut et message
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockClient.AssertExpectations(t)
	})
}

func TestUpdateClient(t *testing.T) {

	t.Run("invalid client data", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Test pour vérifier la validation de l'ID (UUID)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         nil, // ID manquant
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("invalid newsletter data", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Test pour vérifier la validation de la valeur newsletter
		// Aucune interaction avec UpdateClient ici, car la validation échoue avant d'appeler cette méthode
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Valeur incorrecte pour Newsletter
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("successful client update", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mock pour simuler un cas de mise à jour réussie
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(&entities.Client{
			ID: "123e4567-e89b-12d3-a456-426614174000",
		}, nil)

		// Fournir un UUID valide ici
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("client update error", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une erreur lors de la mise à jour du client
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(errors.ErrInternalServer)

		// Fournir un UUID valide mais avec un champ Newsletter manquant (incomplet)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Manque une valeur pour `Newsletter`
		})
		// La validation va échouer avant que `UpdateClient` ne soit appelé, donc nous testons cela
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockClient.AssertNotCalled(t, "UpdateClient") // Vérifier que le mock n'a pas été appelé
	})

	t.Run("client update error", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une erreur lors de la mise à jour du client
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).Return(nil, errors.ErrInternalServer)

		// Fournir un UUID valide mais avec un champ Newsletter manquant (incomplet)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: aws.Bool(true),                                     // Manque une valeur pour `Newsletter`
		})
		// La validation va échouer avant que `UpdateClient` ne soit appelé, donc nous testons cela
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockClient.AssertNotCalled(t, "UpdateClient") // Vérifier que le mock n'a pas été appelé
	})

}

func TestGetClient(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		mockService := new(DomainUserService)
		// Simuler une entrée incorrecte avec un ID manquant
		dtoClient := &transfert.Client{
			ID: nil, // Manque l'ID
		}

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("client not found", func(t *testing.T) {
		mockService := new(DomainUserService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du service qui ne trouve pas le client
		mockService.On("GetClient", dtoClient).Return(nil, errors_domain_user.ErrClientNotFound)

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrClientNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockService := new(DomainUserService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler une erreur interne inattendue
		mockService.On("GetClient", dtoClient).Return(nil, errors.ErrInternalServer)

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockService.AssertExpectations(t)
	})

	t.Run("successful client retrieval", func(t *testing.T) {
		mockService := new(DomainUserService)
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
		mockService := new(DomainUserService)

		// Cas où l'ID du client est manquant
		dtoClient := &transfert.Client{ID: nil}

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("should return 404 if client not found", func(t *testing.T) {
		// Créer un mock pour le service client
		mockService := new(DomainUserService)

		// Cas où le client n'est pas trouvé
		clientID := "123e4567-e89b-12d3-a456-426614174000"
		dtoClient := &transfert.Client{ID: &clientID}

		// Configurer le mock pour renvoyer une erreur client non trouvé
		mockService.On("DeleteClient", dtoClient).Return(errors_domain_user.ErrClientNotFound)

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrClientNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 if internal server error occurs", func(t *testing.T) {

		// Créer un mock pour le service client
		mockService := new(DomainUserService)
		// Cas où une erreur interne survient lors de la suppression
		clientID := "123e4567-e89b-12d3-a456-426614174000"
		dtoClient := &transfert.Client{ID: &clientID}

		// Configurer le mock pour renvoyer une erreur interne
		mockService.On("DeleteClient", dtoClient).Return(errors.ErrInternalServer)

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockService.AssertExpectations(t)
	})

	t.Run("should return 204 if client is deleted successfully", func(t *testing.T) {
		// Créer un mock pour le service client
		mockService := new(DomainUserService)
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
