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
	Config    *Config
	from      string
	expeditor string
}

func (s *Service) From() string {
	return s.from
}

func (s *Service) Expeditor() string {
	return s.expeditor
}

func (s *Service) Send(mail *Mail) error {
	if mail == nil {
		return errors.New("nil mail to send")
	}

	msg, to, err := mail.Prepare(s)
	if err != nil {
		return err
	}

	logger.Info("Sending mail to: ", to)
	return smtp.SendMail(s.Config.Host+":"+s.Config.Port, s.Config.Auth, s.from, to, msg)
}
