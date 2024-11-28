package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
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
		service, _, _, _ := setup()
		require.NotNil(t, service)

		result, err := service.RegisterClient(nil, nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("client already exists", func(t *testing.T) {
		service, mockRepo, _, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("existing@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(&entities.Credential{}, nil)

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential creation error", func(t *testing.T) {
		service, mockRepo, _, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, errors_domain_user.ErrCredentialNotFound)
		mockRepo.On("CreateCredential", dtoCredential).Return(nil, errors.ErrInternalServer)

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "common.internal_error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("client creation error", func(t *testing.T) {
		service, mockRepo, _, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, errors_domain_user.ErrCredentialNotFound)
		mockRepo.On("CreateCredential", dtoCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", dtoClient).Return(nil, errors.ErrInternalServer)

		client, err := service.RegisterClient(dtoCredential, dtoClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "common.internal_error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("client update error", func(t *testing.T) {
		service, mockRepo, _, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, errors_domain_user.ErrCredentialNotFound)
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", inputClient).Return(expectedClient, nil)
		mockRepo.On("UpdateClient", expectedClient).Return(errors.ErrInternalServer)

		client, err := service.RegisterClient(inputCredential, inputClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "common.internal_error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential update error", func(t *testing.T) {
		service, mockRepo, _, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, errors_domain_user.ErrCredentialNotFound)
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateClient", inputClient).Return(expectedClient, nil)
		mockRepo.On("UpdateClient", expectedClient).Return(nil)
		mockRepo.On("UpdateCredential", expectedCredential).Return(errors.ErrInternalServer)

		client, err := service.RegisterClient(inputCredential, inputClient)
		assert.Nil(t, client)
		assert.EqualError(t, err, "common.internal_error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful client and credential creation", func(t *testing.T) {
		service, mockRepo, mockMailer, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, errors_domain_user.ErrCredentialNotFound)
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
	t.Run("no dto", func(t *testing.T) {
		service, _, _, _ := setup()

		// Appel du service avec un DTO nil
		client, err := service.UpdateClient(nil)

		// Vérifier que l'erreur est bien une erreur "No DTO"
		assert.EqualError(t, err, errors.ErrNoDto.Error())
		assert.Nil(t, client)
	})

	t.Run("client not found", func(t *testing.T) {
		service, mockRepo, _, _ := setup()

		// Simuler un client non trouvé dans la base de données
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(nil, errors_domain_user.ErrClientNotFound)

		// Appel du service avec un client introuvable
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("invalid-id")})

		// Vérifier que l'erreur est bien une erreur "Client not found"
		assert.Error(t, err)
		assert.Nil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		service, mockRepo, _, mockPerms := setup()

		// Simuler un client valide
		mockClient := &entities.Client{ID: "42debee6-2063-4566-baf1-37a7bdd139ff"}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler que la méthode CanUpdate retourne false pour ce client
		mockPerms.On("CanUpdate", mockClient, mock.Anything).Return(false)

		// Appel du service avec un client valide mais sans autorisation
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est bien une erreur "Unauthorized"
		assert.EqualError(t, err, errors.ErrUnauthorized.Error())
		assert.Nil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("update client success", func(t *testing.T) {
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
		service, mockRepo, _, mockPerms := setup()

		// Simuler un client valide
		mockClient := &entities.Client{ID: clientID}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler que la méthode CanUpdate retourne true pour ce client
		mockPerms.On("CanUpdate", mockClient, mock.Anything).Return(true)

		// Simuler une mise à jour réussie du client
		mockRepo.On("UpdateClient", mockClient).Return(nil)

		// Appel du service avec un client valide
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est nulle, ce qui signifie que la mise à jour a réussi
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("update client failure", func(t *testing.T) {
		service, mockRepo, _, mockPerms := setup()
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"

		// Simuler un client valide
		mockClient := &entities.Client{ID: clientID}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler que la méthode CanUpdate retourne true
		mockPerms.On("CanUpdate", mockClient, mock.Anything).Return(true)

		// Simuler une erreur lors de la mise à jour du client
		mockRepo.On("UpdateClient", mockClient).Return(errors.ErrInternalServer)

		// Appel du service avec un client valide mais une mise à jour échouée
		client, err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est bien celle retournée par le mock lors de la mise à jour
		assert.EqualError(t, err, "common.internal_error")
		assert.Nil(t, client)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

func TestGetClient(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		service, mockRepo, _, mockPerms := setup()

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

		// Simuler la méthode CanRead pour le contrôle des permissions
		mockPerms.On("CanRead", expectedClient, mock.Anything).Return(true)

		// Appeler la méthode du service
		client, err := service.GetClient(dummyClientDTO)

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Equal(t, expectedClient.ID, client.ID)

		// Vérifier les attentes sur le mock
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("error nil dto", func(t *testing.T) {
		service, _, _, _ := setup()

		// Appeler la méthode du service avec un DTO nil
		client, err := service.GetClient(nil)

		// Vérifier que l'erreur est retournée
		require.Error(t, err)
		require.Nil(t, client)
		assert.EqualError(t, err, errors.ErrNoDto.Error())
	})

	t.Run("cant read client", func(t *testing.T) {
		service, mockRepo, _, mockPerms := setup()

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

		// Simuler la méthode CanRead pour le contrôle des permissions
		mockPerms.On("CanRead", expectedClient, mock.Anything).Return(false)

		// Appeler la méthode du service
		client, err := service.GetClient(dummyClientDTO)

		// Assertions
		require.Error(t, err)
		require.Nil(t, client)

		// Vérifier les attentes sur le mock
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("client_not_found", func(t *testing.T) {
		service, mockRepo, _, _ := setup()

		// Simuler un client DTO valide
		dummyClientDTO := &transfert.Client{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		// Simuler la réponse du repository qui ne trouve pas le client
		mockRepo.On("ReadClient", dummyClientDTO).Return(nil, errors_domain_user.ErrClientNotFound)

		// Appeler la méthode du service
		client, err := service.GetClient(dummyClientDTO)

		// Vérifier que l'erreur est retournée
		require.Error(t, err)
		require.Nil(t, client)
		assert.Error(t, err)

		// Vérifier les attentes sur le mock
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteClient(t *testing.T) {
	t.Run("should return error if dtoClient is nil", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		mockPermission := new(PermissionMock)
		mockGame := new(GameRepositoryMock)
		service := services.User(mockPermission, mockRepo, mockGame, nil)

		err := service.DeleteClient(nil)
		assert.EqualError(t, err, errors.ErrNoDto.Error())
	})

	t.Run("should return error if client not found", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		mockPermission := new(PermissionMock)
		mockGame := new(GameRepositoryMock)
		service := services.User(mockPermission, mockRepo, mockGame, nil)

		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler la lecture du client
		mockRepo.On("ReadClient", dtoClient).Return(nil, errors_domain_user.ErrClientNotFound)

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier que l'erreur est bien celle attendue
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if client cannot be deleted", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		mockPermission := new(PermissionMock)
		mockGame := new(GameRepositoryMock)
		service := services.User(mockPermission, mockRepo, mockGame, nil)

		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler la lecture du client
		mockRepo.On("ReadClient", dtoClient).Return(&entities.Client{ID: *clientID}, nil)
		// Simuler la permission de suppression
		mockPermission.On("CanDelete", mock.AnythingOfType("*entities.Client")).Return(false)

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier que l'erreur est bien celle attendue
		assert.EqualError(t, err, errors.ErrUnauthorized.Error())
		mockRepo.AssertExpectations(t)
		mockPermission.AssertExpectations(t)
	})

	t.Run("should delete client successfully", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		mockPermission := new(PermissionMock)
		mockGame := new(GameRepositoryMock)
		service := services.User(mockPermission, mockRepo, mockGame, nil)
		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler la lecture du client
		mockRepo.On("ReadClient", dtoClient).Return(&entities.Client{ID: *clientID}, nil)
		// Simuler la permission de suppression
		mockPermission.On("CanDelete", mock.AnythingOfType("*entities.Client")).Return(true)
		// Simuler la suppression réussie du client
		mockRepo.On("DeleteClient", dtoClient).Return(nil)

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier qu'il n'y a pas d'erreur
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPermission.AssertExpectations(t)
	})

	t.Run("should return error if repository delete fails", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		mockPermission := new(PermissionMock)
		mockGame := new(GameRepositoryMock)
		service := services.User(mockPermission, mockRepo, mockGame, nil)

		// Client DTO avec un ID valide
		clientID := aws.String("123e4567-e89b-12d3-a456-426614174000")
		dtoClient := &transfert.Client{ID: clientID}

		// Simuler la lecture du client
		mockRepo.On("ReadClient", dtoClient).Return(&entities.Client{ID: *clientID}, nil)
		// Simuler la permission de suppression
		mockPermission.On("CanDelete", mock.AnythingOfType("*entities.Client")).Return(true)
		// Simuler une erreur lors de la suppression du client
		mockRepo.On("DeleteClient", dtoClient).Return(errors.ErrInternalServer)

		// Appel du service pour supprimer le client
		err := service.DeleteClient(dtoClient)

		// Vérifier que l'erreur est bien celle attendue
		assert.EqualError(t, err, "common.internal_error")
		mockRepo.AssertExpectations(t)
		mockPermission.AssertExpectations(t)
	})
}
