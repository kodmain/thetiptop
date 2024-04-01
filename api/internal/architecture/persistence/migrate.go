package persistence

import (
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
)

func Migrate(models ...any) {
	logger.Warn("migrating models to databases")
	for _, model := range models {
		for database := range databases {
			if err := databases[database].AutoMigrate(model); err != nil {
				application.PANIC <- err
			}
			logger.Warnf("migrated model %T to database %s", model, database)
		}
	}
}
