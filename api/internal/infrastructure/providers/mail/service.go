package mail

import (
	"errors"
	"net/smtp"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

type ServiceInterface interface {
	Send(*Mail) error
	From() string
	Expeditor() string
}

type Service struct {
	Config *Config
}

func (s *Service) From() string {
	if s.Config == nil {
		return ""
	}

	return s.Config.From
}

func (s *Service) Expeditor() string {
	if s.Config == nil {
		return ""
	}

	return s.Config.Expeditor
}

func (s *Service) Send(mail *Mail) error {
	if mail == nil {
		return errors.New("nil mail to send")
	}

	if s.Config == nil {
		return errors.New("nil config")
	}

	msg, to, err := mail.Prepare(s)
	if err != nil {
		return err
	}

	logger.Info("Sending mail to: ", to)
	return smtp.SendMail(s.Config.Host+":"+s.Config.Port, s.Config.Auth, s.From(), to, msg)
}
