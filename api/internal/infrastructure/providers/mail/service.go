package mail

import (
	"errors"
	"net/smtp"
)

// Service Gère la configuration du service de messagerie et l'envoi des mails.
//
// Fields:
// - Host: string L'hôte SMTP.
// - Port: string Le port SMTP.
// - Username: string Le nom d'utilisateur pour l'authentification SMTP.
// - Password: string Le mot de passe pour l'authentification SMTP.
// - Auth: smtp.Auth L'authentification SMTP.
//
// Methods:
// - New: Initialise le service avec la configuration donnée.
// - Send: Envoie un e-mail.
type Service struct {
	Host      string
	Port      string
	Username  string
	Password  string
	From      string
	Expeditor string
	Auth      smtp.Auth
}

var instance *Service

// New Initialise le service de messagerie avec la configuration donnée.
//
// Parameters:
// - cfg: *Service La configuration du service de messagerie.
//
// Returns:
// - error: Une erreur si l'initialisation échoue.
func New(cfg *Service) error {
	if cfg == nil {
		return errors.New("mail configuration is nil")
	}

	if cfg.Host == "" {
		return errors.New("mail host is empty")
	}

	if cfg.Port == "" {
		return errors.New("mail port is empty")
	}

	if cfg.From == "" {
		return errors.New("mail from is empty")
	}

	instance = cfg
	instance.Auth = smtp.PlainAuth("", instance.Username, instance.Password, instance.Host)

	return nil
}

// Send Envoie un e-mail en utilisant la configuration du service.
//
// Parameters:
// - mail: *Mail L'e-mail à envoyer.
//
// Returns:
// - error: Une erreur si l'envoi échoue.
func Send(mail *Mail) error {
	if instance == nil {
		return errors.New("mail instance is nil")
	}

	if !mail.IsValid() {
		return errors.New("invalid mail")
	}

	msg, to, err := mail.Prepare()
	if err != nil {
		return err
	}

	return smtp.SendMail(instance.Host+":"+instance.Port, instance.Auth, instance.From, to, msg)
}