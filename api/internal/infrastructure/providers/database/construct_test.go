package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConfigNew(t *testing.T) {

	db := database.Get("noinstance")
	assert.Nil(t, db)

	err := database.New(nil)
	assert.Error(t, err)

	databases := map[string]*database.Config{
		"memory": {
			Protocol: database.SQLite,
			DBname:   CONF_MEMORY,
			Logger:   true,
		},
		"mysql": {
			Protocol: database.MySQL,
			Host:     CONF_HOST,
			Port:     CONF_MYSQL_PORT,
			User:     CONF_USER,
			Password: CONF_PASSWORD,
			DBname:   CONF_DBNAME,
		},
		"postgres": {
			Protocol: database.PostgreSQL,
			Host:     CONF_HOST,
			Port:     CONF_PG_PORT,
			User:     CONF_USER,
			Password: CONF_PASSWORD,
			DBname:   CONF_DBNAME,
		},
		"unknown": {
			Protocol: "unknown",
		},
		"empty": nil,
	}

	for key, cfg := range databases {
		err = database.New(map[string]*database.Config{key: cfg})
		if key == "empty" || key == "unknown" {
			assert.Error(t, err)
		}
	}

	for key, cfg := range databases {
		err = database.New(map[string]*database.Config{key: cfg})
		assert.Error(t, err)
	}

	db = database.Get("memory")
	assert.NotNil(t, db)
	db = database.Get("empty")
	assert.Nil(t, db)
	db = database.Get()
	assert.Nil(t, db)
	db = database.Get("notexist")
	assert.Nil(t, db)
}

func TestConfigFrom(t *testing.T) {
	database.FromDB(nil)

	// Création du mock SQL
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Création de l'instance Gorm avec le mock
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Création de l'instance de Database avec le mock
	dbInstance, err := database.FromDB(gormDB)
	require.NoError(t, err)
	require.NotNil(t, dbInstance)
}
