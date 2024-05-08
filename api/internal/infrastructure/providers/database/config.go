package database

import (
	"fmt"
	"strconv"
	"strings"
)

type Options map[string]string

const (
	MySQL      string = "mysql"
	PostgreSQL string = "postgres"
	SQLite     string = "sqlite"
)

type Config struct {
	Protocol string `yaml:"protocol"` // 'mysql', 'postgres', 'sqlite'
	Host     string `yaml:"host"`     // 'localhost', '127.0.0.1', ou vide pour SQLite
	Port     string `yaml:"port"`     // '3306', '5432', ou vide pour SQLite
	User     string `yaml:"user"`     // Nom d'utilisateur, vide pour SQLite
	Password string `yaml:"password"` // Mot de passe, vide pour SQLite
	DBname   string `yaml:"dbname"`   // Nom de la base de données, chemin vers le fichier SQLite ou ':memory:'

	// Paramètres spécifiques pour la connexion à la base de données
	// Par exemple, pour MySQL & PostgreSQL, ça pourrait être 'sslmode=disable', etc.
	Options Options
}

func (cfg *Config) ToDSN() (string, error) {
	switch cfg.Protocol {
	case MySQL:
		// Format DSN MySQL : user:password@tcp(host:port)/dbname?options
		return strings.TrimSpace(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, formatOptions(cfg.Options))), nil
	case PostgreSQL:
		// Format DSN PostgreSQL : h=myh port=myport user=myuser dbname=mydb options
		return strings.TrimSpace(fmt.Sprintf("host=%s port=%s user=%s dbname=%s %s", cfg.Host, cfg.Port, cfg.User, cfg.DBname, formatOptions(cfg.Options))), nil
	case SQLite:
		// SQLite utilise simplement le chemin du fichier ou ':memory:'
		return strings.TrimSpace(cfg.DBname), nil
	default:
		return "", fmt.Errorf("protocole inconnu")
	}
}

func (cfg *Config) Validate() error {
	switch cfg.Protocol {
	case MySQL:
		return cfg.validateMySQL()
	case PostgreSQL:
		return cfg.validatePostgreSQL()
	case SQLite:
		return cfg.validateSQLite()
	default:
		return fmt.Errorf("protocole inconnu")
	}
}

func (cfg *Config) common() error {
	if cfg.Host == "" {
		return fmt.Errorf("missing host for " + cfg.Protocol)
	}

	if !isValidPort(cfg.Port) {
		return fmt.Errorf("missing port for " + cfg.Protocol)
	}

	if cfg.Password == "" {
		return fmt.Errorf("missing password for " + cfg.Protocol)
	}

	if cfg.User == "" || cfg.User == "root" {
		return fmt.Errorf("wrong user " + cfg.Protocol)
	}

	return nil
}

func (cfg *Config) validateMySQL() error {
	return cfg.common()
}

func (cfg *Config) validatePostgreSQL() error {
	return cfg.common()
}

func (cfg *Config) validateSQLite() error {
	if cfg.DBname == "" {
		return fmt.Errorf("le nom de la base de données SQLite est requis")
	}
	return nil
}

func isValidPort(port string) bool {
	p, err := strconv.Atoi(port)
	return err == nil && p > 0 && p <= 65535
}

func formatOptions(options map[string]string) string {
	if len(options) == 0 {
		return ""
	}

	var opts []string
	for k, v := range options {
		opts = append(opts, fmt.Sprintf("%s=%s", k, v))
	}

	return "?" + strings.Join(opts, "&")
}
