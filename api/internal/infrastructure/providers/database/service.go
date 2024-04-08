package database

import (
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
	for key, cfg := range *databases {
		if cfg == nil {
			return fmt.Errorf("database configuration is required")
		}

		if _, ok := instances[key]; ok {
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

		/*
			import glogger "gorm.io/gorm/logger"
			newLogger := glogger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer (stdout)
				glogger.Config{
					SlowThreshold: time.Second,  // Seuil de temps lent pour les requêtes
					LogLevel:      glogger.Info, // LogLevel Info pour voir toutes les requêtes
					Colorful:      true,         // Activer les couleurs
				},
			)
		*/

		instances[key], err = gorm.Open(dial, &gorm.Config{
			PrepareStmt: true,
			Logger:      glogger.Discard,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func Get(name string) *gorm.DB {
	if _, ok := instances[name]; !ok {
		application.PANIC <- fmt.Errorf("database not found")
	}

	return instances[name]
}
