package database_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseServices(t *testing.T) {
	err := database.New(nil)
	assert.Error(t, err)

	databases := database.Databases{
		"sql": &database.Database{
			Protocol: database.MySQL,
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "password",
			DBname:   "mydb",
		},
		"test": &database.Database{
			Protocol: database.SQLite,
			DBname:   ":memory:",
		},
		"mysql": &database.Database{
			Protocol: database.MySQL,
			Host:     "localhost",
			Port:     "3306",
			User:     "user",
			Password: "password",
			DBname:   "mydb",
		},
		"postgres": &database.Database{
			Protocol: database.PostgreSQL,
			Host:     "localhost",
			Port:     "3306",
			User:     "user",
			Password: "password",
			DBname:   "mydb",
		},
		"unknow": &database.Database{
			Protocol: "unknow",
		},
		"empty": nil,
	}

	err = database.New(&databases)
	assert.Error(t, err)

	err = database.New(&databases)
	assert.Error(t, err)

	db := database.Get("test")
	assert.NotNil(t, db)
	db = database.Get("oki")
	assert.Nil(t, db)

}
