package database

import (
	"errors"
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	glogger "gorm.io/gorm/logger"
)

var instances map[string]ServiceInterface = make(map[string]ServiceInterface)

func New(databases map[string]*Config) error {
	if databases == nil {
		return fmt.Errorf("database configuration is required")
	}

	errs := make([]error, 0)

	for key, cfg := range databases {
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

		db, err := gorm.Open(dial, &gorm.Config{
			PrepareStmt: true,
			Logger:      glogger.Discard,
		})

		if err != nil {
			errs = append(errs, err)
		}

		instances[key] = &Service{
			Config: cfg,
			db:     db,
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
