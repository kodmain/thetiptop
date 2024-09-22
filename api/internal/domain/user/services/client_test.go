package services_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClientRegister(t *testing.T) {
	// Variables communes
	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2061-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	sidClient := idClient.String()

	inputClient := &transfert.Client{
		CGU: aws.Bool(true),
	}

	expectedClient := &entities.Client{
		ID:  idClient.String(),
		CGU: aws.Bool(true),
		Validations: []*entities.Validation{
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("666666").Pointer(),
				Type:      0,
				Validated: false,
				ClientID:  &sidClient,
			},
		},
	}

	inputCredential := &transfert.Credential{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	expectedCredential := &entities.Credential{
		ID:       idClient.String(),
		Email:    inputCredential.Email,
		Password: aws.String("$2a$10$wO5PfDAGp6w2ubKp0vEdXeUe2HlfOv5iRJ3C3MVR0vJhscD0G.NKS"), // hashed password
		//ClientID: &sidClient,
	}

	t.Run("nil input", func(t *testing.T) {
		service, _, _ := setup()
		require.NotNil(t, service)

		result, err := service.RegisterClient(nil, nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrNoDto, err.Error())
	})

	t.Run("client already exists", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("existing@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(&entities.Credential{}, nil)

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, errors.ErrClientAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential creation error", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", dtoCredential).Return(nil, fmt.Errorf("error creating credential"))

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "error creating credential")
		mockRepo.AssertExpectations(t)
	})

	t.Run("client creation error", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", dtoCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", dtoClient).Return(nil, fmt.Errorf("error creating client"))

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "error creating client")
		mockRepo.AssertExpectations(t)
	})

	t.Run("client update error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", inputClient).Return(expectedClient, nil)
		mockRepo.On("UpdateClient", expectedClient).Return(fmt.Errorf("error updating client"))

		client, err := service.RegisterClient(inputCredential, inputClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "error updating client")
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential update error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", inputClient).Return(expectedClient, nil)
		mockRepo.On("UpdateClient", expectedClient).Return(nil)
		mockRepo.On("UpdateCredential", expectedCredential).Return(fmt.Errorf("error updating credential"))

		client, err := service.RegisterClient(inputCredential, inputClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "error updating credential")
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful client and credential creation", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", inputClient).Return(expectedClient, nil)
		mockRepo.On("UpdateClient", expectedClient).Return(nil)
		mockRepo.On("UpdateCredential", expectedCredential).Return(nil)

		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		client, err := service.RegisterClient(inputCredential, inputClient)
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, sidClient, client.ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateClient(t *testing.T) {
	t.Run("client not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un client non trouvé dans la base de données
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel du service avec un client introuvable
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("invalid-id")})

		// Vérifier que l'erreur est bien une erreur "Client not found"
		assert.EqualError(t, err, errors.ErrClientNotFound)
		assert.Nil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})

	t.Run("update client success", func(t *testing.T) {
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
		service, mockRepo, _ := setup()

		// Simuler un client valide
		mockClient := &entities.Client{ID: clientID}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler une mise à jour réussie du client
		mockRepo.On("UpdateClient", mockClient).Return(nil)

		// Appel du service avec un client valide
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est nulle, ce qui signifie que la mise à jour a réussi
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})

	t.Run("update client failure", func(t *testing.T) {
		service, mockRepo, _ := setup()
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"

		// Simuler un client valide
		mockClient := &entities.Client{ID: clientID}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler une erreur lors de la mise à jour du client
		mockRepo.On("UpdateClient", mockClient).Return(fmt.Errorf("update error"))

		// Appel du service avec un client valide mais une mise à jour échouée
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est bien celle retournée par le mock lors de la mise à jour
		assert.EqualError(t, err, "update error")
		assert.Nil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})
}

func TestGetClient(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		service, mockRepo, _ := setup()
		// Simuler un client DTO valide
		dummyClientDTO := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}
		expectedClient := &entities.Client{
			ID:  "42debee6-2063-4566-baf1-37a7bdd139ff",
			CGU: aws.Bool(true),
		}

		// Simuler la réponse du repository
		mockRepo.On("ReadClient", dummyClientDTO).Return(expectedClient, nil)

		// Appeler la méthode du service
		client, err := service.GetClient(dummyClientDTO)

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Equal(t, expectedClient.ID, client.ID)

		// Vérifier les attentes sur le mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("error_nil_dto", func(t *testing.T) {
		service, _, _ := setup()
		// Appeler la méthode du service avec un DTO nil
		client, err := service.GetClient(nil)

		// Vérifier que l'erreur est retournée
		require.Error(t, err)
		require.Nil(t, client)
		assert.EqualError(t, err, errors.ErrNoDto)
	})

	t.Run("client_not_found", func(t *testing.T) {
		service, mockRepo, _ := setup()
		// Simuler un client DTO valide
		dummyClientDTO := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du repository qui ne trouve pas le client
		mockRepo.On("ReadClient", dummyClientDTO).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appeler la méthode du service
		client, err := service.GetClient(dummyClientDTO)

		// Vérifier que l'erreur est retournée
		require.Error(t, err)
		require.Nil(t, client)
		assert.EqualError(t, err, errors.ErrClientNotFound)

		// Vérifier les attentes sur le mock
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteClient(t *testing.T) {
	ctx := new(fiber.Ctx)

	t.Run("should return error if dtoClient is nil", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		service := services.User(ctx, mockRepo, nil)

		err := service.DeleteClient(nil)
		assert.EqualError(t, err, errors.ErrNoDto)
	})

	t.Run("should return error if client ID is nil", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		service := services.User(ctx, mockRepo, nil)

		dtoClient := &transfert.Client{ID: nil}
		err := service.DeleteClient(dtoClient)
		assert.EqualError(t, err, errors.ErrNoDto)
	})

	t.Run("should delete client successfully", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		service := services.User(ctx, mockRepo, nil)
		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler la suppression réussie du client
		mockRepo.On("DeleteClient", dtoClient).Return(nil)

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier qu'il n'y a pas d'erreur
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if repository delete fails", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		service := services.User(ctx, mockRepo, nil)

		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler une erreur lors de la suppression du client
		mockRepo.On("DeleteClient", dtoClient).Return(fmt.Errorf("delete failed"))

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier que l'erreur est bien celle attendue
		assert.EqualError(t, err, "delete failed")
		mockRepo.AssertExpectations(t)
	})
}
