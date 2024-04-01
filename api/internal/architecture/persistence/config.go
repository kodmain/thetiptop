package persistence

import (
	"fmt"
	"strconv"
	"strings"
)

type Protocol int

const (
	MySQL Protocol = iota
	PostgreSQL
	SQLite
)

type Config struct {
	Name string // Nom de la base de données

	Protocol Protocol // 'mysql', 'postgres', 'sqlite'
	Host     string   // 'localhost', '127.0.0.1', ou vide pour SQLite
	Port     string   // '3306', '5432', ou vide pour SQLite
	User     string   // Nom d'utilisateur, vide pour SQLite
	Password string   // Mot de passe, vide pour SQLite
	DBName   string   // Nom de la base de données, chemin vers le fichier SQLite ou ':memory:'

	// Paramètres spécifiques pour la connexion à la base de données
	// Par exemple, pour MySQL & PostgreSQL, ça pourrait être 'sslmode=disable', etc.
	Options map[string]string
}

func (cfg *Config) ToDSN() (string, error) {
	switch cfg.Protocol {
	case MySQL:
		// Format DSN MySQL : user:password@tcp(host:port)/dbname?options
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, formatOptions(cfg.Options)), nil
	case PostgreSQL:
		// Format DSN PostgreSQL : host=myhost port=myport user=myuser dbname=mydb options
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s %s", cfg.Host, cfg.Port, cfg.User, cfg.DBName, formatOptions(cfg.Options)), nil
	case SQLite:
		// SQLite utilise simplement le chemin du fichier ou ':memory:'
		return cfg.DBName, nil
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

func (cfg *Config) validateMySQL() error {
	if cfg.Host == "" {
		return fmt.Errorf("l'hôte MySQL est requis")
	}
	if !isValidPort(cfg.Port) {
		return fmt.Errorf("le port MySQL est invalide ou manquant")
	}
	if cfg.User == "" {
		return fmt.Errorf("l'utilisateur MySQL est requis")
	}

	return nil
}

func (cfg *Config) validatePostgreSQL() error {
	if cfg.Host == "" {
		return fmt.Errorf("l'hôte PostgreSQL est requis")
	}
	if !isValidPort(cfg.Port) {
		return fmt.Errorf("le port PostgreSQL est invalide ou manquant")
	}
	if cfg.User == "" {
		return fmt.Errorf("l'utilisateur PostgreSQL est requis")
	}

	return nil
}

func (cfg *Config) validateSQLite() error {
	if cfg.DBName == "" {
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
