package mail

import (
	"net/smtp"
)

type Config struct {
	Host      string
	Port      string
	Username  string
	Password  string
	From      string
	Expeditor string
	Auth      smtp.Auth
}
