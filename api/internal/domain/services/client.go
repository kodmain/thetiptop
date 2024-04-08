package services

import (
	"errors"

	"github.com/kodmain/thetiptop/api/internal/application/dto"
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
		repo: repositories.NewClientRepository("default"),
		mail: mail.NewTemplate("signup"),
	}

	return instance
}

func (s *ClientService) SignUp(dto *dto.Client) error {
	_, err := s.repo.Read(dto)
	if err == nil {
		return errors.New("client already exists")
	}

	client, err := s.repo.Create(dto)
	if err != nil {
		return err
	}

	text, html, err := s.mail.Inject(map[string]interface{}{
		"AppName": "The Tip Top",
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
