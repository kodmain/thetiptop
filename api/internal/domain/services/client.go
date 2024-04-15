package services

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	interfaces "github.com/kodmain/thetiptop/api/internal/domain/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/repositories"
)

var instance *ClientService

type ClientService struct {
	repo interfaces.ClientRepository
	mail *mail.Template
}

func Client() *ClientService {
	if instance != nil {
		return instance
	}

	instance = &ClientService{
		repo: repositories.NewClientRepository(config.DEFAULT_DB_NAME),
		mail: mail.NewTemplate("signup"),
	}

	return instance
}

func (s *ClientService) SignUp(obj *transfert.Client) error {
	_, err := s.repo.Read(obj)
	if err == nil {
		return err
	}

	client, err := s.repo.Create(obj)
	if err != nil {
		return err
	}

	text, html, err := s.mail.Inject(mail.Data{
		"AppName": "ThéTipTop",
		"Url":     "https://thetiptop.com",
	})

	if err != nil {
		return err
	}

	m := &mail.Mail{
		To:      []string{client.Email},
		Subject: "Welcome to The Tip Top",
		Text:    text,
		Html:    html,
	}

	return mail.Send(m)
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
