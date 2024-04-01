package persistence

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/architecture/events"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var databases map[string]*gorm.DB = make(map[string]*gorm.DB)

func New(cfgs ...Config) error {
	for _, cfg := range cfgs {
		if _, ok := databases[cfg.Name]; ok {
			return fmt.Errorf("database already exists")
		}

		if err := cfg.Validate(); err != nil {
			return err
		}

		dsn, err := cfg.ToDSN()
		if err != nil {
			return fmt.Errorf("database name is required")
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
			return fmt.Errorf("unknown protocol")
		}

		databases[cfg.Name], err = gorm.Open(dial, &gorm.Config{
			PrepareStmt: true,
		})

		if err != nil {
			return err
		}

		events.Notify(events.MIGRATE, cfg.Name)
	}

	return nil
}

func Get(name string) *gorm.DB {
	if _, ok := databases[name]; !ok {
		application.PANIC <- fmt.Errorf("database not found")
	}

	return databases[name]
}
