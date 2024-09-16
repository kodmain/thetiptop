package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
		//ClientID: &clientID,
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
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(expectedCredential, nil)

		// Simuler un client valide
		mockRepo.On("ReadClient", mock.AnythingOfType("*transfert.Client")).
			Return(expectedClient, nil)

		// Appel du service avec un credential valide
		result, err := service.UserAuth(inputCredential)

		// Vérification des résultats
		require.NoError(t, err)
		require.NotNil(t, result)

		// Vérifier que les attentes sur le mock sont satisfaites
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

		// Mock pour ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)

		// Mock pour UpdateValidation
		mockRepo.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérification qu'il n'y a pas d'erreur et que la validation a été retournée
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Vérification que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("error client not found", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Mock pour `ReadValidation`
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).
			Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérification qu'une erreur est présente et que le résultat est nul
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérification que toutes les attentes du mock sont respectées
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

		// Simuler une réponse d'erreur pour ReadValidation
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(nil, fmt.Errorf(errors.ErrValidationNotFound)).Once()

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il y a une erreur et que le résultat est nul
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("update fail", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil).Once()

		// Simuler une erreur lors de la mise à jour de la validation
		mockRepo.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(fmt.Errorf("update error")).Once()

		// Appel de la méthode PasswordValidation
		result, err := service.PasswordValidation(&transfert.Validation{}, &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("newpassword123"),
		})

		// Vérifier qu'il y a une erreur et que le résultat est nul
		assert.Error(t, err)
		assert.Nil(t, result)

		// Vérifier que toutes les attentes du mock sont respectées
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("validation expired", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

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

		// Simuler une réponse réussie de la méthode ReadValidation
		mockValidation := &entities.Validation{
			ClientID:  aws.String("valid-client-id"),
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Validated: false,
		}
		mockRepo.On("ReadValidation", mock.AnythingOfType("*transfert.Validation")).Return(mockValidation, nil)
		// Appel de la méthode PasswordValidation
		result, err := service.MailValidation(&transfert.Validation{}, &transfert.Credential{
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

func TestValidationRecover(t *testing.T) {

	t.Run("credential fail", func(t *testing.T) {
		service, mockRepo, _ := setup()

		err := service.ValidationRecover(nil, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier que l'erreur correspond à "Client not found"
		assert.EqualError(t, err, errors.ErrNoDto)
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler que le credential n'est pas trouvé
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		// Vérifier que l'erreur correspond à "Client not found"
		assert.EqualError(t, err, errors.ErrUserNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success client", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		luhn := token.Generate(6)

		// Simuler que le credential n'est pas trouvé
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(&entities.Credential{
				Email: aws.String("test@example.com"), // L'email est vide, car on cherche à simuler un credential non trouvé
			}, nil)

		mockRepo.On("ReadUser", mock.AnythingOfType("*transfert.User")).
			Return(&entities.Client{}, nil, nil)

		mockRepo.On("CreateValidation", mock.AnythingOfType("*transfert.Validation")).
			Return(&entities.Validation{
				Token: &luhn,
			}, nil)

		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success employee", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		luhn := token.Generate(6)

		// Simuler que le credential n'est pas trouvé
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(&entities.Credential{
				Email: aws.String("test@example.com"), // L'email est vide, car on cherche à simuler un credential non trouvé
			}, nil)

		mockRepo.On("ReadUser", mock.AnythingOfType("*transfert.User")).
			Return(nil, &entities.Employee{}, nil)

		mockRepo.On("CreateValidation", mock.AnythingOfType("*transfert.Validation")).
			Return(&entities.Validation{
				Token: &luhn,
			}, nil)

		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail no client or employee", func(t *testing.T) {
		service, mockRepo, _ := setup()

		// Simuler que le credential n'est pas trouvé
		mockRepo.On("ReadCredential", mock.AnythingOfType("*transfert.Credential")).
			Return(&entities.Credential{
				Email: aws.String("test@example.com"), // L'email est vide, car on cherche à simuler un credential non trouvé
			}, nil)

		mockRepo.On("ReadUser", mock.AnythingOfType("*transfert.User")).
			Return(nil, nil, fmt.Errorf(errors.ErrUserNotFound))

		err := service.ValidationRecover(&transfert.Validation{}, &transfert.Credential{
			Email: aws.String("test@example.com"),
		})

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
