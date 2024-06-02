package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	glogger "gorm.io/gorm/logger"
)

var (
	instances map[string]*Database = make(map[string]*Database)
	mutex     sync.RWMutex
)

func FromDB(db *gorm.DB) (*Database, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &Database{Engine: db}, nil
}

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

		mutex.Lock()
		if _, ok := instances[key]; ok {
			mutex.Unlock()
			errs = append(errs, fmt.Errorf("database already exists"))
			continue
		}
		mutex.Unlock()

		if err := cfg.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}

		dsn := cfg.ToDSN()

		var dial gorm.Dialector
		logger.Warnf("connecting to %s", dsn)
		switch cfg.Protocol {
		case SQLite:
			dial = sqlite.Open(dsn)
		case MySQL:
			dial = mysql.Open(dsn)
		case PostgreSQL:
			dial = postgres.Open(dsn)
		}

		db, err := gorm.Open(dial, &gorm.Config{
			PrepareStmt: true,
			Logger:      glogger.Discard,
		})

		if err != nil {
			errs = append(errs, err)
		}

		mutex.Lock()
		instances[key] = &Database{
			Config: cfg,
			Engine: db,
		}
		mutex.Unlock()
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func Get(names ...string) *Database {
	mutex.RLock()
	defer mutex.RUnlock()

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
