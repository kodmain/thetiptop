package services

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
)

type ClientServiceInterface interface {
	// Sign
	SignUp(obj *transfert.Client) (*entities.Client, error)
	SignIn(obj *transfert.Client) (*entities.Client, error)
	SignValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error)

	// Password
	PasswordRecover(obj *transfert.Client) error
	PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error)
}

type ClientService struct {
	repo repositories.ClientRepositoryInterface
	mail mail.ServiceInterface
}

func Client(repo repositories.ClientRepositoryInterface, mail mail.ServiceInterface) *ClientService {
	return &ClientService{repo, mail}
}

func (s *ClientService) SignUp(obj *transfert.Client) (*entities.Client, error) {
	if obj == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	_, err := s.repo.ReadClient(obj)
	if err == nil {
		return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
	}

	client, err := s.repo.CreateClient(obj)
	if err != nil {
		return nil, err
	}

	go s.sendSignUpMail(client)

	return client, nil
}

func (s *ClientService) SignIn(obj *transfert.Client) (*entities.Client, error) {
	client, err := s.repo.ReadClient(&transfert.Client{
		Email: obj.Email,
	})

	if err != nil || !client.CompareHash(*obj.Password) {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	if validation := client.HasSuccessValidation(entities.MailValidation); validation == nil {
		return nil, fmt.Errorf(errors.ErrClientNotValidate, entities.MailValidation.String())
	}

	return client, nil
}

func (s *ClientService) PasswordRecover(obj *transfert.Client) error {
	query := &transfert.Client{
		Email: obj.Email,
	}

	client, err := s.repo.ReadClient(query)
	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
	}

	if mailValidation := client.HasSuccessValidation(entities.MailValidation); mailValidation == nil {
		return fmt.Errorf(errors.ErrClientNotValidate, entities.MailValidation.String())
	}

	client.Validations = append(client.Validations, &entities.Validation{
		Type: entities.PasswordRecover,
	})

	if err := s.repo.UpdateClient(client); err != nil {
		return err
	}

	go s.sendMailRecover(client)

	return nil
}

func (s *ClientService) PasswordUpdate(obj *transfert.Client) error {
	client, err := s.repo.ReadClient(&transfert.Client{
		Email: obj.Email,
	})

	if err != nil {
		return fmt.Errorf(errors.ErrClientNotFound)
	}

	if passwordValidation := client.HasSuccessValidation(entities.PasswordRecover); passwordValidation == nil {
		return fmt.Errorf(errors.ErrClientNotValidate, entities.PasswordRecover.String())
	}

	password, err := hash.Hash(aws.String(*client.Email+":"+*obj.Password), hash.BCRYPT)
	if err != nil {
		return err
	}

	client.Password = password

	if err := s.repo.UpdateClient(client); err != nil {
		return err
	}

	return nil
}

// sendMail Send a templated email to a client
// This function handles the common logic for sending templated emails to clients.
//
// Parameters:
// - client: *entities.Client The client to send the email to.
// - templateName: string The name of the email template.
// - validationType: entities.ValidationType The type of validation to check.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *ClientService) sendMail(client *entities.Client, templateName string, validationType entities.ValidationType) error {
	tpl := template.NewTemplate(templateName)
	if tpl == nil {
		return fmt.Errorf(errors.ErrTemplateNotFound, templateName)
	}

	validation := client.HasNotExpiredValidation(validationType)
	if validation == nil {
		return fmt.Errorf(errors.ErrValidationNotFound)
	}

	text, html, err := tpl.Inject(template.Data{
		"AppName": env.APP_NAME,
		"Token":   validation.Token.String(),
	})
	if err != nil {
		return err
	}

	subject := ""
	if validationType == entities.PasswordRecover {
		subject = "Récupération de mot de passe"
	} else if validationType == entities.MailValidation {
		subject = "Bienvenue chez The Tip Top"
	}

	m := &mail.Mail{
		To:      []string{*client.Email},
		Subject: subject,
		Text:    text,
		Html:    html,
	}

	for i := 0; i < 3; i++ {
		if err := s.mail.Send(m); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf(errors.ErrMailSendFailed)
}

// sendMailRecover Send a password recovery email to a client
// This function sends a password recovery email to the specified client.
//
// Parameters:
// - client: *entities.Client The client to send the email to.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *ClientService) sendMailRecover(client *entities.Client) error {
	return s.sendMail(client, "recover", entities.PasswordRecover)
}

// sendSignUpMail Send a signup confirmation email to a client
// This function sends a signup confirmation email to the specified client.
//
// Parameters:
// - client: *entities.Client The client to send the email to.
//
// Returns:
// - error: error An error object if an error occurs, nil otherwise.
func (s *ClientService) sendSignUpMail(client *entities.Client) error {
	return s.sendMail(client, "signup", entities.MailValidation)
}
