package database_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabases(t *testing.T) {

	dbconfs := []struct {
		failOrSuccess bool
		protocol      string
		host          string
		port          string
		password      string
		user          string
		dbname        string
		option        map[string]string
	}{
		// MySQL
		{true, database.MySQL, "localhost", "3306", "password", "user", "mydb", database.Options{"sslmode": "disable"}},
		{false, database.MySQL, "", "3306", "password", "root", "mydb", nil},
		{false, database.MySQL, "localhost", "", "password", "root", "mydb", nil},
		{false, database.MySQL, "localhost", "3306", "", "root", "mydb", nil},
		{false, database.MySQL, "localhost", "3306", "password", "", "mydb", nil},
		{false, database.MySQL, "localhost", "3306", "password", "root", "", nil},
		// PostgreSQL
		{true, database.PostgreSQL, "localhost", "5432", "password", "user", "mydb", database.Options{"sslmode": "disable"}},
		{false, database.PostgreSQL, "", "5432", "password", "root", "mydb", nil},
		{false, database.PostgreSQL, "localhost", "", "password", "root", "mydb", nil},
		{false, database.PostgreSQL, "localhost", "5432", "", "root", "mydb", nil},
		{false, database.PostgreSQL, "localhost", "5432", "password", "", "mydb", nil},
		{false, database.PostgreSQL, "localhost", "5432", "password", "root", "", nil},
		// SQLite
		{true, database.SQLite, "", "", "", "", "hello.db", nil},
		{false, database.SQLite, "", "", "", "", "", nil},
		// Unknown
		{false, "unknow", "", "", "", "", "", nil},
	}

	for _, dbconf := range dbconfs {
		db := &database.Database{
			Protocol: dbconf.protocol,
			Host:     dbconf.host,
			Port:     dbconf.port,
			User:     dbconf.user,
			Password: dbconf.password,
			DBname:   dbconf.dbname,
			Options:  dbconf.option,
		}

		dsn, err := db.ToDSN()
		if dbconf.protocol != "unknow" {
			assert.NoError(t, err)
			assert.NotNil(t, dsn)
		}
		assert.Equal(t, dbconf.failOrSuccess, db.Validate() == nil, dsn)
	}
}
