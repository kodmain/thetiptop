package services_test

import (
	"fmt"
	"reflect"
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
		t.Parallel()
		mockClient := new(DomainUserService)
		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("short"),
		}, &transfert.Client{
			Newsletter: aws.Bool(true),
			CGU:        aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "password", "errors.Errors should contain a key 'password'")
		}
		mockClient.AssertExpectations(t)
	})

	t.Run("missing newsletter", func(t *testing.T) {
		t.Parallel()
		mockClient := new(DomainUserService)
		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidPass123!"),
		}, &transfert.Client{
			Newsletter: nil,
			CGU:        aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "newsletter", "errors.Errors should contain a key 'newsletter'")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("invalid cgu", func(t *testing.T) {
		t.Parallel()
		mockClient := new(DomainUserService)
		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidPass123!"),
		}, &transfert.Client{
			Newsletter: aws.Bool(true),
			CGU:        aws.Bool(false),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "cgu", "errors.Errors should contain a key 'cgu'")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("valid password and fields", func(t *testing.T) {
		t.Parallel()
		mockClient := new(DomainUserService)
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(&entities.Client{}, nil)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidPass123!"),
		}, &transfert.Client{
			Newsletter: aws.Bool(true),
			CGU:        aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusCreated, statusCode)
		assert.NotNil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("client already exists", func(t *testing.T) {
		t.Parallel()
		mockClient := new(DomainUserService)
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(nil, errors_domain_user.ErrCredentialAlreadyExists)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidPass123!"),
		}, &transfert.Client{
			Newsletter: aws.Bool(true),
			CGU:        aws.Bool(true),
		})

		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)

		mockClient.AssertExpectations(t)
	})

	t.Run("server error during registration", func(t *testing.T) {
		t.Parallel()
		mockClient := new(DomainUserService)
		mockClient.On("RegisterClient", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Client")).Return(nil, errors.ErrInternalServer)

		statusCode, response := services.RegisterClient(mockClient, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidPass123!"),
		}, &transfert.Client{
			Newsletter: aws.Bool(true),
			CGU:        aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		err, ok := response.(errors.ErrorInterface)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Equal(t, errors.ErrInternalServer, err)
		}

		mockClient.AssertExpectations(t)
	})
}

// TestMailValidation tests the MailValidation and ExportClient services
// This test suite checks various scenarios for email validation and client export
//
// Parameters:
// - t: *testing.T test framework
//
// Returns:
// - none
func TestMailValidation(t *testing.T) {
	// Chargement de la configuration
	config.Load(aws.String("../../../config.test.yml"))

	luhn := token.Generate(6)
	email := "valid.email@example.com"

	// Cas de validation réussie
	t.Run("successful validation", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)

		// Génération d'un ID aléatoire pour la validation
		id, err := uuid.NewRandom()
		assert.NoError(t, err)

		// Simulation de la réponse du mock pour une validation réussie
		mockClient.On("MailValidation", mock.Anything, mock.Anything).
			Return(&entities.Validation{ID: id.String()}, nil)

		// Appel de la fonction à tester
		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: luhn.PointerString()},
			&transfert.Credential{Email: aws.String(email)},
		)

		// Vérifications
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
		assert.IsType(t, &entities.Validation{}, response)

		// Vérification des attentes du mock
		mockClient.AssertExpectations(t)
	})

	// Cas d'un token avec une syntaxe invalide
	t.Run("invalid token syntax", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)

		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: aws.String("invalidToken")},
			&transfert.Credential{Email: aws.String(emailSyntaxFail)},
		)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// La réponse est un objet JSON contenant une clé "token"
		// Exemple: {"token": "validator.is_not_number"}
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "token", "Key 'token' should exist in errors.Errors")
		}

		mockClient.AssertExpectations(t)
	})

	// Cas d'un email avec une syntaxe invalide
	t.Run("invalid email syntax", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)

		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: luhn.PointerString()},
			&transfert.Credential{Email: aws.String(emailSyntaxFail)},
		)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// La réponse est un objet JSON contenant une clé "email"
		// Exemple: {"email": "validator.is_not_email"}
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "email", "Key 'email' should exist in errors.Errors")
		}

		mockClient.AssertExpectations(t)
	})

	// Cas où la validation n'a pas été trouvée
	t.Run("validation not found", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		mockClient.On("MailValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationNotFound)

		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: luhn.PointerString()},
			&transfert.Credential{Email: aws.String(email)},
		)

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationNotFound, response)

		mockClient.AssertExpectations(t)
	})

	// Cas où la validation a expiré (nommé ici "validation already validated" dans ton code)
	// Tu peux renommer en "validation expired" si tu le souhaites
	t.Run("validation already validated", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		mockClient.On("MailValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationExpired)

		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: luhn.PointerString()},
			&transfert.Credential{Email: aws.String(email)},
		)

		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationExpired, response)

		mockClient.AssertExpectations(t)
	})

	// Cas où la validation a déjà été effectuée
	t.Run("validation already validated", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		mockClient.On("MailValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationAlreadyValidated)

		statusCode, response := services.MailValidation(
			mockClient,
			&transfert.Validation{Token: luhn.PointerString()},
			&transfert.Credential{Email: aws.String(email)},
		)

		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationAlreadyValidated, response)

		mockClient.AssertExpectations(t)
	})

	// Tests liés à l'export du client
	t.Run("should return 200 and client data on success", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)

		// Simuler des données de client exportées
		expectedClients := &entities.ClientData{
			Credential:  &entities.Credential{ID: "credential-id", Email: aws.String("email@example.com")},
			Client:      &entities.Client{ID: "client-id", CGU: aws.Bool(true)},
			Validations: []*entities.Validation{{ID: "validation-id", Validated: false}},
			Tickets:     []*entitiesGame.Ticket{{ID: "ticket-id"}},
		}

		mockService.On("ExportClient").Return(expectedClients, nil)

		statusCode, response := services.ExportClient(mockService)

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedClients, response)

		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on internal error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("ExportClient").Return(nil, errors.ErrInternalServer)

		statusCode, response := services.ExportClient(mockService)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)

		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on unauthorized", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("ExportClient").Return(nil, errors.ErrUnauthorized)

		statusCode, response := services.ExportClient(mockService)

		assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		assert.Equal(t, errors.ErrUnauthorized, response)

		mockService.AssertExpectations(t)
	})

	t.Run("should return error code and error on client not found", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("ExportClient").Return(nil, errors.ErrNotFound)

		statusCode, response := services.ExportClient(mockService)

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors.ErrNotFound, response)

		mockService.AssertExpectations(t)
	})
}

func TestValidationRecover(t *testing.T) {
	// On charge la configuration
	config.Load(aws.String("../../../config.test.yml"))

	validationType := "password_recovery"
	email := "valid.email@example.com"
	emailSyntaxFail := "invalid-email"

	t.Run("invalid email", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(
			mockClient,
			&transfert.Credential{Email: aws.String(emailSyntaxFail)},
			&transfert.Validation{Type: aws.String(validationType)},
		)

		// Vérifie le Status
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// Vérifie qu'on a bien une map d'erreurs
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			// On peut ensuite vérifier une clé "email" par exemple
			assert.Contains(t, errorsMap, "email")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("missing validation type", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Cas d'une erreur de validation de l'email
		statusCode, response := services.ValidationRecover(
			mockClient,
			&transfert.Credential{Email: aws.String(email)},
			&transfert.Validation{Type: nil},
		)

		fmt.Println("type", reflect.TypeOf(response))
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errorsMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errorsMap, "type")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("successful recovery", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Mock pour un cas de récupération réussi
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).
			Return(nil)

		statusCode, response := services.ValidationRecover(
			mockClient,
			&transfert.Credential{
				Email: aws.String(email),
			},
			&transfert.Validation{
				Type: aws.String(validationType),
			},
		)

		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("client not found", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Mock pour simuler un client non trouvé
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).
			Return(errors_domain_user.ErrUserNotFound)

		statusCode, response := services.ValidationRecover(
			mockClient,
			&transfert.Credential{
				Email: aws.String(email),
			},
			&transfert.Validation{
				Type: aws.String(validationType),
			},
		)

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrUserNotFound, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("other error", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Mock pour simuler une erreur interne
		mockClient.On("ValidationRecover", mock.AnythingOfType("*transfert.Validation"), mock.AnythingOfType("*transfert.Credential")).
			Return(errors.ErrInternalServer)

		statusCode, response := services.ValidationRecover(
			mockClient,
			&transfert.Credential{
				Email: aws.String(email),
			},
			&transfert.Validation{
				Type: aws.String(validationType),
			},
		)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.NotEmpty(t, errObj.Error(), "The error message should not be empty")
		}

		mockClient.AssertExpectations(t)
	})
}
func TestUpdateClient(t *testing.T) {
	t.Run("invalid client data", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)

		// Test pour vérifier la validation de l'ID (UUID)
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         nil, // ID manquant
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// On s'attend à recevoir un type errors.Errors si ta validation renvoie un "map"
		// ou *errors.Error si c'est un objet unique. Adapte selon ce que renvoie ta logique.
		errMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMap, "id", "Key 'id' should exist in validation errors")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("invalid newsletter data", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Test pour vérifier la validation de la valeur Newsletter
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Valeur incorrecte/absente
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// Même logique ici : adapter le cast si besoin
		errMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMap, "newsletter", "Key 'newsletter' should exist in validation errors")
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("successful client update", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Mock pour simuler un cas de mise à jour réussie
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).
			Return(&entities.Client{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			}, nil)

		// Fournir un UUID valide ici
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: aws.Bool(true),
		})
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		// Ici, on s'attend à recevoir un *entities.Client.
		clientData, ok := response.(*entities.Client)
		assert.True(t, ok, "response should be of type *entities.Client")
		if ok {
			assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", clientData.ID)
		}

		mockClient.AssertExpectations(t)
	})

	t.Run("client update validation error", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Dans ce scénario, on veut que la validation échoue avant d'appeler UpdateClient,
		// donc on fournit une donnée invalide, par exemple Newsletter = nil
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"), // UUID valide
			Newsletter: nil,                                                // Validation doit échouer
		})

		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// Vérifie qu'on obtient des erreurs de validation
		errMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMap, "newsletter", "Key 'newsletter' should exist in validation errors")
		}

		// Vérifie que le mock n'a pas été appelé, puisque la validation échoue
		mockClient.AssertNotCalled(t, "UpdateClient")
	})

	t.Run("client update internal error", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Simuler une erreur lors de la mise à jour du client
		mockClient.On("UpdateClient", mock.AnythingOfType("*transfert.Client")).
			Return(nil, errors.ErrInternalServer)

		// Ici, on passe des données valides pour que la validation réussisse
		// Ensuite, le service va renvoyer une erreur
		statusCode, response := services.UpdateClient(mockClient, &transfert.Client{
			ID:         aws.String("123e4567-e89b-12d3-a456-426614174000"),
			Newsletter: aws.Bool(true),
		})

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)

		// Selon ton code, `UpdateClient` peut renvoyer un *errors.Error ou errors.Errors
		// Adapte en conséquence
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Equal(t, errors.ErrInternalServer, errObj)
		}

		// Ici, on s'attend à ce que le mock ait été appelé, car la validation a passé
		mockClient.AssertCalled(t, "UpdateClient", mock.AnythingOfType("*transfert.Client"))
		mockClient.AssertExpectations(t)
	})
}

func TestGetClient(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		// Simuler une entrée incorrecte avec un ID manquant
		dtoClient := &transfert.Client{
			ID: nil, // Manque l'ID
		}

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		// Selon la structure de ton code, si tu renvoies un *errors.Error :
		mapErrors, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.NotEmpty(t, mapErrors, "errors.Errors should not be empty")
		}

		mockService.AssertExpectations(t)
	})

	t.Run("client not found", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du service qui ne trouve pas le client
		mockService.On("GetClient", dtoClient).
			Return(nil, errors_domain_user.ErrClientNotFound)

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrClientNotFound, response)

		mockService.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		// Simuler une entrée valide
		dtoClient := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler une erreur interne inattendue
		mockService.On("GetClient", dtoClient).
			Return(nil, errors.ErrInternalServer)

		// Appel de la méthode
		statusCode, response := services.GetClient(mockService, dtoClient)

		// Vérifications
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)

		// Idem, si c'est un *errors.Error :
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Equal(t, errors.ErrInternalServer, errObj)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("successful client retrieval", func(t *testing.T) {
		t.Parallel()

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

		mockService.On("GetClient", dtoClient).
			Return(expectedClient, nil)

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
		t.Parallel()

		mockService := new(DomainUserService)

		// Cas où l'ID du client est manquant
		dtoClient := &transfert.Client{ID: nil}

		// Appel de la fonction DeleteClient
		statusCode, response := services.DeleteClient(mockService, dtoClient)

		// Vérifier le résultat
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		errMaps, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errMaps)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 if client not found", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("should return 204 if client is deleted successfully", func(t *testing.T) {
		t.Parallel()

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
