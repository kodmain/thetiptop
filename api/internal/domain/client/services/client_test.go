package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ClientRepositoryMock struct {
	mock.Mock
}

func (m *ClientRepositoryMock) CreateClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *ClientRepositoryMock) ReadClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *ClientRepositoryMock) UpdateClient(client *entities.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientRepositoryMock) DeleteClient(client *transfert.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientRepositoryMock) CreateValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *ClientRepositoryMock) ReadValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *ClientRepositoryMock) UpdateValidation(validation *entities.Validation) error {
	args := m.Called(validation)

	if args.Get(0) == nil {
		validation.ID = uuid.New().String()
		validation.Token = token.NewLuhn("666666").Pointer()
		return nil
	}

	return args.Error(0)
}

func (m *ClientRepositoryMock) DeleteValidation(validation *transfert.Validation) error {
	args := m.Called(validation)
	return args.Error(0)
}

func (m *ClientRepositoryMock) CreateCredential(credential *transfert.Credential) (*entities.Credential, error) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credential), args.Error(1)
}

func (m *ClientRepositoryMock) ReadCredential(credential *transfert.Credential) (*entities.Credential, error) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credential), args.Error(1)
}

func (m *ClientRepositoryMock) UpdateCredential(credential *entities.Credential) error {
	args := m.Called(credential)
	if args.Get(0) == nil {
		credential.ID = uuid.New().String()
		return nil
	}
	return args.Error(0)
}

func (m *ClientRepositoryMock) DeleteCredential(credential *transfert.Credential) error {
	args := m.Called(credential)
	return args.Error(0)
}

type MailServiceMock struct {
	mock.Mock
}

func (m *MailServiceMock) Send(mail *mail.Mail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func (m *MailServiceMock) From() string {
	args := m.Called()
	return args.String(0)
}

func (m *MailServiceMock) Expeditor() string {
	args := m.Called()
	return args.String(0)
}

func setup() (*services.ClientService, *ClientRepositoryMock, *MailServiceMock) {
	mockRepository := new(ClientRepositoryMock)
	mockMailer := new(MailServiceMock)
	service := services.Client(mockRepository, mockMailer)

	return service, mockRepository, mockMailer
}
func TestUserRegister(t *testing.T) {
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
		ClientID: &sidClient,
	}

	t.Run("nil input", func(t *testing.T) {
		service, _, _ := setup()
		require.NotNil(t, service)

		result, err := service.UserRegister(nil, nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrNoDto, err.Error())
	})

	t.Run("client already exists", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("existing@example.com")}
		dtoClient := &transfert.Client{}

		mockRepo.On("ReadCredential", dtoCredential).Return(&entities.Credential{}, nil)

		client, err := service.UserRegister(dtoCredential, dtoClient)
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

		client, err := service.UserRegister(dtoCredential, dtoClient)
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

		client, err := service.UserRegister(dtoCredential, dtoClient)
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

		client, err := service.UserRegister(inputCredential, inputClient)
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

		client, err := service.UserRegister(inputCredential, inputClient)
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

		client, err := service.UserRegister(inputCredential, inputClient)
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.Equal(t, sidClient, client.ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserAuth(t *testing.T) {
	// Variables communes
	email := aws.String("test@example.com")
	password := aws.String("password123")
	failpassword := aws.String("password1234")
	hashedPassword, err := hash.Hash(aws.String(*email+":"+*password), hash.BCRYPT)
	require.NoError(t, err)
	clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
	credentialID := "42debee6-2063-4566-baf1-37a7bdd139f0"

	inputCredential := &transfert.Credential{
		ID:       &credentialID,
		Email:    email,
		Password: password,
	}

	// Simuler un credential valide
	expectedCredential := &entities.Credential{
		Email:    email,
		Password: hashedPassword,
		ClientID: &clientID,
	}

	// Simuler un client valide
	expectedClient := &entities.Client{
		ID:  clientID,
		CGU: aws.Bool(true),
		Validations: []*entities.Validation{
			{
				Type:      entities.MailValidation,
				Validated: true,
			},
		},
	}

	t.Run("credential not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential non trouvé, peu importe les valeurs spécifiques des champs
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(nil, fmt.Errorf(errors.ErrCredentialNotFound))

		// Appeler UserAuth avec un credential dont l'ID est nil (pour simuler un credential non trouvé)
		client, err := service.UserAuth(&transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("wrongpassword"),
			ID:       nil, // L'ID est nil, car on cherche à simuler un credential non trouvé
		})

		// Vérification que le client est nul et que l'erreur correspond à "ErrCredentialNotFound"
		assert.Nil(t, client)
		assert.EqualError(t, err, errors.ErrCredentialNotFound)

		// Vérifier que les attentes sur le mock sont satisfaites
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential valide mais un mot de passe incorrect
		expectedCredential := &entities.Credential{
			Email:    email,
			Password: aws.String("hashedCorrectPassword"), // Le mot de passe correct haché
		}

		// Simuler un credential valide
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(expectedCredential, nil)

		// Appeler le service avec un mot de passe incorrect
		client, err := service.UserAuth(&transfert.Credential{
			Email:    email,
			Password: aws.String("wrongpassword"),
		})

		// Vérification que le client est nul et que l'erreur concerne un mot de passe incorrect
		assert.Nil(t, client)
		assert.EqualError(t, err, errors.ErrCredentialNotFound) // Assurez-vous que l'erreur est appropriée
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential hash fail", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential valide mais un échec de hachage
		expectedCredential := &entities.Credential{
			Email:    email,
			Password: aws.String("hashedCorrectPassword"), // Mot de passe haché correct
		}

		// Le mock retourne un credential valide
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(expectedCredential, nil)

		// Appel du service avec un mot de passe incorrect
		client, err := service.UserAuth(&transfert.Credential{
			Email:    email,
			Password: failpassword, // Mot de passe incorrect
		})

		// Vérifier que le client est nul et que l'erreur concerne un mot de passe incorrect
		assert.Nil(t, client)
		assert.EqualError(t, err, errors.ErrCredentialNotFound) // Utiliser l'erreur correcte
		mockRepo.AssertExpectations(t)
	})

	t.Run("client not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un appel `ReadCredential` qui retourne le credential attendu
		mockRepo.On("ReadCredential", mock.MatchedBy(func(cred *transfert.Credential) bool {
			return cred.Email != nil && *cred.Email == *email
		})).Return(expectedCredential, nil)

		// Simuler que le client n'est pas trouvé
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel du service avec un credential valide
		result, err := service.UserAuth(inputCredential)

		// Vérification des résultats
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientNotFound, err.Error())

		// Vérifier que les attentes sur le mock sont satisfaites
		mockRepo.AssertExpectations(t)
	})

	t.Run("client found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un appel `ReadCredential` qui retourne le credential attendu
		mockRepo.On("ReadCredential", mock.MatchedBy(func(cred *transfert.Credential) bool {
			return cred.Email != nil && *cred.Email == *email
		})).Return(expectedCredential, nil)

		// Simuler un client valide
		mockRepo.On("ReadClient", mock.MatchedBy(func(client *transfert.Client) bool {
			return client.ID != nil && *client.ID == clientID // Comparer en déréférençant client.ID
		})).Return(expectedClient, nil)

		// Appel du service avec un credential valide
		result, err := service.UserAuth(inputCredential)

		// Vérification des résultats
		require.NoError(t, err)
		require.NotNil(t, result)

		// Vérifier que les attentes sur le mock sont satisfaites
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
		err := service.UpdateClient(&transfert.Client{ID: aws.String("invalid-id")})

		// Vérifier que l'erreur est bien une erreur "Client not found"
		assert.EqualError(t, err, errors.ErrClientNotFound)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})

	t.Run("update client success", func(t *testing.T) {

		service, mockRepo, _ := setup()

		// Simuler un client valide
		mockClient := &entities.Client{ID: "42debee6-2063-4566-baf1-37a7bdd139ff"}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler une mise à jour réussie du client
		mockRepo.On("UpdateClient", mockClient).Return(nil)

		// Appel du service avec un client valide
		err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est nulle, ce qui signifie que la mise à jour a réussi
		assert.NoError(t, err)

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})

	t.Run("update client failure", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un client valide
		mockClient := &entities.Client{ID: "42debee6-2063-4566-baf1-37a7bdd139ff"}

		// Le mock retourne un client valide pour l'appel à ReadClient
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler une erreur lors de la mise à jour du client
		mockRepo.On("UpdateClient", mockClient).Return(fmt.Errorf("update error"))

		// Appel du service avec un client valide mais une mise à jour échouée
		err := service.UpdateClient(&transfert.Client{ID: aws.String("valid-id")})

		// Vérifier que l'erreur est bien celle retournée par le mock lors de la mise à jour
		assert.EqualError(t, err, "update error")

		// Vérifier que les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
	})
}
func TestValidationRecover(t *testing.T) {

	t.Run("credential not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler que le credential n'est pas trouvé
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier que l'erreur correspond à "Client not found"
		assert.EqualError(t, err, errors.ErrClientNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("client not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential valide
		mockCredential := &entities.Credential{ClientID: aws.String("valid-client-id")}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(mockCredential, nil)

		// Simuler que le client n'est pas trouvé
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier que l'erreur correspond à "Client not found"
		assert.EqualError(t, err, errors.ErrClientNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation update success", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential et un client valides avec un ClientID synchronisé
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
		mockCredential := &entities.Credential{ClientID: aws.String(clientID)}
		mockClient := &entities.Client{ID: clientID}
		mockValidation := &entities.Validation{ClientID: aws.String(clientID)}

		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(mockCredential, nil)
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler la création de la validation
		mockRepo.On("CreateValidation", mock.AnythingOfType("*transfert.Validation")).
			Return(mockValidation, nil) // Retourne une validation et pas d'erreur

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier qu'il n'y a pas d'erreur
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation create fail", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential et un client valides avec un ClientID synchronisé
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
		mockCredential := &entities.Credential{ClientID: aws.String(clientID)}
		mockClient := &entities.Client{ID: clientID}

		// Générer l'objet Validation attendu
		expectedValidation := &transfert.Validation{
			ClientID: aws.String(clientID),
			// Remplir les autres champs selon ce que CreateValidation ferait
		}

		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(mockCredential, nil)
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Utiliser mock.MatchedBy pour matcher dynamiquement l'objet Validation lors de la création
		mockRepo.On("CreateValidation", mock.MatchedBy(func(v *transfert.Validation) bool {
			// Vérifier que le ClientID correspond bien
			return v.ClientID != nil && *v.ClientID == *expectedValidation.ClientID
		})).Return(nil, fmt.Errorf("validation create error")) // Simuler l'erreur lors de la création

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier que l'erreur correspond à l'erreur de création de la validation
		assert.EqualError(t, err, "validation create error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation update success", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un credential et un client valides avec un ClientID synchronisé
		clientID := "42debee6-2063-4566-baf1-37a7bdd139ff"
		mockCredential := &entities.Credential{ClientID: aws.String(clientID)}
		mockClient := &entities.Client{ID: clientID}

		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(mockCredential, nil)
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(mockClient, nil)

		// Simuler la création de la validation
		mockRepo.On("CreateValidation", mock.AnythingOfType("*transfert.Validation")).
			Return(nil, fmt.Errorf("update")) // Retourne une validation et pas d'erreur

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier qu'il n'y a pas d'erreur
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

}

func TestPasswordUpdate(t *testing.T) {
	t.Run("TestPasswordUpdate_Success", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un Credential existant
		mockCredential := &entities.Credential{Email: aws.String("test@example.com"), Password: aws.String("old-password")}
		newPassword := "new-password"

		// Simuler la lecture réussie du credential
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)
		// Simuler le hachage réussi du nouveau mot de passe
		// Simuler la mise à jour réussie du credential
		mockRepo.On("UpdateCredential", mockCredential).Return(nil)

		// Appel de la méthode PasswordUpdate
		err := service.PasswordUpdate(&transfert.Credential{Email: mockCredential.Email, Password: aws.String(newPassword)})

		// Vérifier qu'il n'y a pas d'erreur
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestPasswordUpdate_Fail", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler un Credential existant
		mockCredential := &entities.Credential{Email: aws.String("test@example.com"), Password: aws.String("old-password")}
		newPassword := "new-password"

		// Simuler la lecture réussie du credential
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)
		// Simuler le hachage réussi du nouveau mot de passe
		// Simuler la mise à jour réussie du credential
		mockRepo.On("UpdateCredential", mockCredential).Return(fmt.Errorf("update error"))

		// Appel de la méthode PasswordUpdate
		err := service.PasswordUpdate(&transfert.Credential{Email: mockCredential.Email, Password: aws.String(newPassword)})

		// Vérifier qu'il n'y a pas d'erreur
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestPasswordUpdate_ClientNotFound", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler une erreur de type ErrClientNotFound
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel de la méthode PasswordUpdate
		err := service.PasswordUpdate(&transfert.Credential{Email: aws.String("test@example.com"), Password: aws.String("new-password")})

		// Vérifier que l'erreur correspond bien à ErrClientNotFound
		assert.EqualError(t, err, errors.ErrClientNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("no dto", func(t *testing.T) {
		service, _, _ := setup()
		err := service.PasswordUpdate(nil)
		assert.EqualError(t, err, errors.ErrNoDto)
	})

	t.Run("no dto", func(t *testing.T) {
		service, _, _ := setup()
		err := service.PasswordUpdate(nil)
		assert.EqualError(t, err, errors.ErrNoDto)
	})
}

func TestPasswordValidation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		mockRepo.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("error client not found", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("no dto", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(nil, nil)

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("validation not found", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("update fail", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		mockRepo.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(fmt.Errorf("update error"))
		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("validation expired", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("validation expired", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("validation expired", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadCredential
		mockCredential := &entities.Credential{
			ClientID: aws.String("valid-client-id"),
			Email:    aws.String("test@example.com"),
			Password: aws.String("hashed-password"),
		}
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).Return(mockCredential, nil)

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Validated: true,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})
}
