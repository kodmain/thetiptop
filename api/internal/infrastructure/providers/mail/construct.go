package mail

import (
	"errors"
	"net/smtp"
)

var instances map[string]ServiceInterface = make(map[string]ServiceInterface)

// New Initialise le service de messagerie avec la configuration donnée.
//
// Parameters:
// - cfg: *Config La configuration du service de messagerie.
//
// Returns:
// - error: Une erreur si l'initialisation échoue.
func New(mailers map[string]*Config) error {
	if mailers == nil {
		return errors.New("mail configuration is required")
	}

	errs := make([]error, 0)

	for name, cfg := range mailers {
		if cfg == nil {
			errs = append(errs, errors.New("mail "+name+" config is nil"))
			continue
		}

		if cfg.Host == "" {
			errs = append(errs, errors.New("mail host is empty"))
		}

		if cfg.Port == "" {
			errs = append(errs, errors.New("mail port is empty"))
		}

		if cfg.From == "" {
			errs = append(errs, errors.New("mail from is empty"))
		}

		if cfg.Username != "" && cfg.Password != "" {
			cfg.Auth = smtp.Auth(smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host))
		}

		instances[name] = &Service{
			Config: cfg,
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func Get(names ...string) ServiceInterface {
	if len(instances) == 0 {
		return nil
	}

	var name string
	if len(names) != 1 {
		name = "default"
	} else {
		name = names[0]
	}

	return instances[name]
}
