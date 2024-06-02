package database_test

import (
	"fmt"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
)

const (
	CONF_HOST       = "localhost"
	CONF_MYSQL_PORT = "3306"
	CONF_PG_PORT    = "5432"
	CONF_USER       = "user"
	CONF_ROOT       = "root"
	CONF_PASSWORD   = "password"
	CONF_DBNAME     = "mydb"
	CONF_MEMORY     = ":memory:"
	CONF_FILE       = "file.db"
	CONF_EMPTY      = ""
)

func TestDatabasesConfiguration(t *testing.T) {

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
		{true, database.MySQL, CONF_HOST, CONF_MYSQL_PORT, CONF_PASSWORD, CONF_USER, CONF_DBNAME, database.Options{"sslmode": "disable"}},
		{false, database.MySQL, CONF_EMPTY, CONF_MYSQL_PORT, CONF_PASSWORD, CONF_USER, CONF_DBNAME, nil},
		{false, database.MySQL, CONF_HOST, CONF_EMPTY, CONF_PASSWORD, CONF_USER, CONF_DBNAME, nil},
		{false, database.MySQL, CONF_HOST, CONF_MYSQL_PORT, CONF_EMPTY, CONF_USER, CONF_DBNAME, nil},
		{false, database.MySQL, CONF_HOST, CONF_MYSQL_PORT, CONF_PASSWORD, CONF_EMPTY, CONF_DBNAME, nil},
		{false, database.MySQL, CONF_HOST, CONF_MYSQL_PORT, CONF_PASSWORD, CONF_ROOT, CONF_EMPTY, nil},
		// PostgreSQL
		{true, database.PostgreSQL, CONF_HOST, CONF_PG_PORT, CONF_PASSWORD, CONF_USER, CONF_DBNAME, database.Options{"sslmode": "disable"}},
		{false, database.PostgreSQL, CONF_EMPTY, CONF_PG_PORT, CONF_PASSWORD, CONF_USER, CONF_DBNAME, nil},
		{false, database.PostgreSQL, CONF_HOST, CONF_EMPTY, CONF_PASSWORD, CONF_USER, CONF_DBNAME, nil},
		{false, database.PostgreSQL, CONF_HOST, CONF_PG_PORT, CONF_EMPTY, CONF_USER, CONF_DBNAME, nil},
		{false, database.PostgreSQL, CONF_HOST, CONF_PG_PORT, CONF_PASSWORD, CONF_EMPTY, CONF_DBNAME, nil},
		{false, database.PostgreSQL, CONF_HOST, CONF_PG_PORT, CONF_PASSWORD, CONF_ROOT, CONF_EMPTY, nil},
		// SQLite
		{true, database.SQLite, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_FILE, nil},
		{false, database.SQLite, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, nil},
		// Unknown
		{false, "unknown", CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, CONF_EMPTY, nil},
	}

	for idx, dbconf := range dbconfs {
		db := &database.Config{
			Protocol: dbconf.protocol,
			Host:     dbconf.host,
			Port:     dbconf.port,
			User:     dbconf.user,
			Password: dbconf.password,
			DBname:   dbconf.dbname,
			Options:  dbconf.option,
		}

		dsn := db.ToDSN()
		if dbconf.protocol != "unknown" {
			assert.NotNil(t, dsn)
		}

		fmt.Println("test:", idx)
		assert.Equal(t, dbconf.failOrSuccess, db.Validate() == nil, dsn)
	}
}
