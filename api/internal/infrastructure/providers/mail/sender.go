package mail

import "net/smtp"

type Sender struct {
}

func (s *Sender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

type SenderInterface interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}
