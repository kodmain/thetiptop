package services

import (
	"fmt"
	"time"

	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
)

type ClientServiceInterface interface {
	ValidationMail(obj *transfert.Validation) (*entities.Validation, error)
	SignUp(obj *transfert.Client) (*entities.Client, error)
	SignIn(obj *transfert.Client) (*entities.Client, error)
}

type ClientService struct {
	repo repositories.ClientRepositoryInterface
	mail mail.ServiceInterface
}

func Client(repo repositories.ClientRepositoryInterface, mail mail.ServiceInterface) *ClientService {
	return &ClientService{repo, mail}
}

func (s *ClientService) ValidationMail(obj *transfert.Validation) (*entities.Validation, error) {
	validation, err := s.repo.ReadValidation(obj)

	if err != nil {
		return nil, fmt.Errorf(errors.ErrValidationNotFound)
	}

	if validation.Validated {
		return nil, fmt.Errorf(errors.ErrValidationAlreadyValidated)
	}

	if validation.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf(errors.ErrValidationExpired)
	}

	validation.Validated = true

	if s.repo.UpdateValidation(validation) != nil {
		return nil, err
	}

	return validation, nil
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

func (s *ClientService) sendSignUpMail(client *entities.Client) error {
	tpl := template.NewTemplate("signup")

	if len(client.Validations) == 0 {
		return fmt.Errorf(errors.ErrValidationNotFound)
	}

	text, html, err := tpl.Inject(template.Data{
		"AppName": env.APP_NAME,
		"Url":     env.HOSTNAME + ":" + fmt.Sprintf("%d", *env.PORT_HTTP) + "/validation/" + client.ID + "/" + client.Validations[0].Token.String(),
	})

	if err != nil {
		return err
	}

	m := &mail.Mail{
		To:      []string{*client.Email},
		Subject: "Welcome to The Tip Top",
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

func (s *ClientService) SignIn(obj *transfert.Client) (*entities.Client, error) {
	client, err := s.repo.ReadClient(&transfert.Client{
		Email: obj.Email,
	})

	if err != nil || !client.CompareHash(*obj.Password) {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	validation := client.Validations.Has(entities.Mail)
	if validation == nil || !validation.Validated {
		return nil, fmt.Errorf(errors.ErrClientNotValidated)
	}

	return client, nil
}
