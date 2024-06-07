package services

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
)

type ClientService struct {
	repo repositories.ClientRepository
	mail mail.ServiceInterface
}

func Client(repo repositories.ClientRepository, mail mail.ServiceInterface) *ClientService {
	return &ClientService{repo, mail}
}

func (s *ClientService) SignUp(obj *transfert.Client) (*entities.Client, error) {
	if obj == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	_, err := s.repo.Read(obj)
	if err == nil {
		return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
	}

	client, err := s.repo.Create(obj)
	if err != nil {
		return nil, err
	}

	go s.sendMail(client)

	return client, nil
}

func (s *ClientService) sendMail(client *entities.Client) {
	tpl := template.NewTemplate("signup")

	text, html, err := tpl.Inject(template.Data{
		"AppName": env.APP_NAME,
		"Url":     env.HOSTNAME,
	})

	if !logger.Error(err) {
		m := &mail.Mail{
			To:      []string{client.Email},
			Subject: "Welcome to The Tip Top",
			Text:    text,
			Html:    html,
		}

		// if failed retry 3 times
		logger.Error(s.mail.Send(m))
	}
}

func (s *ClientService) SignIn(obj *transfert.Client) (*entities.Client, error) {
	client, err := s.repo.Read(&transfert.Client{
		Email: obj.Email,
	})

	if err != nil || !client.CompareHash(obj.Password) {
		return nil, fmt.Errorf(errors.ErrClientNotFound)
	}

	return client, nil
}
