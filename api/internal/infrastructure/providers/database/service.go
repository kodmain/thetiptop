package database

import (
	"errors"
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var instances map[string]*gorm.DB = make(map[string]*gorm.DB)

func New(databases *Databases) error {
	if databases == nil {
		return fmt.Errorf("database configuration is required")
	}

	errs := make([]error, 0)

	for key, cfg := range *databases {
		if cfg == nil {
			errs = append(errs, fmt.Errorf("database configuration is required"))
			continue
		}

		if _, ok := instances[key]; ok {
			errs = append(errs, fmt.Errorf("database already exists"))
			continue
		}

		if err := cfg.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}

		dsn, err := cfg.ToDSN()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		var dial gorm.Dialector
		logger.Warnf("connecting to %s", dsn)
		switch cfg.Protocol {
		case SQLite:
			dial = sqlite.Open(dsn)
		case MySQL:
			dial = mysql.Open(dsn)
		case PostgreSQL:
			dial = postgres.Open(dsn)
		default:
			errs = append(errs, fmt.Errorf("unknown protocol"))
		}

		instances[key], err = gorm.Open(dial, &gorm.Config{
			PrepareStmt: true,
			Logger:      glogger.Discard,
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func Get(name string) *gorm.DB {
	if _, ok := instances[name]; !ok {
		application.PANIC <- fmt.Errorf("database not found")
	}

	return instances[name]
}
