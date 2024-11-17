package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/kodmain/thetiptop/api/internal/application/hook"
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
		if err := initializeDatabase(key, cfg, &errs); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func initializeDatabase(key string, cfg *Config, errs *[]error) error {
	if cfg == nil {
		return fmt.Errorf("database configuration for %s is required", key)
	}

	if isInstanceExists(key) {
		return nil
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	dsn := cfg.ToDSN()
	dial, err := getGormDialector(cfg.Protocol, dsn)
	if err != nil {
		return err
	}

	logger.Warnf("connecting to %s", dsn)
	gcfg := buildGormConfig(cfg.Logger)

	db, err := gorm.Open(dial, gcfg)
	if err != nil {
		return err
	}

	hook.Call(hook.EventOnDBInit)
	saveDatabaseInstance(key, cfg, db)

	return nil
}

func isInstanceExists(key string) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	_, exists := instances[key]
	return exists
}

func getGormDialector(protocol string, dsn string) (gorm.Dialector, error) {
	switch protocol {
	case SQLite:
		return sqlite.Open(dsn), nil
	case MySQL:
		return mysql.Open(dsn), nil
	case PostgreSQL:
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

func buildGormConfig(enableLogger bool) *gorm.Config {
	config := &gorm.Config{
		PrepareStmt: true,
		Logger:      glogger.Discard,
	}

	if enableLogger {
		config.Logger = glogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), glogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  glogger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
			ParameterizedQueries:      false,
		})
	}

	return config
}

func saveDatabaseInstance(key string, cfg *Config, db *gorm.DB) {
	mutex.Lock()
	defer mutex.Unlock()

	instances[key] = &Database{
		Config: cfg,
		Engine: db,
	}
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
