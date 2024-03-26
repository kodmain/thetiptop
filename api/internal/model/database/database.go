package database

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	ERROR_INVALID_CONF = "invalid database configuration for: %s -> %s"
	ERROR_INVALID_DSN  = "invalid database URI connection string"

	TYPE_SQLITE    = "sqlite"
	TYPE_MYSQL     = "mysql"
	TYPE_POSTGRES  = "postgres"
	TYPE_SQLSERVER = "sqlserver"

	REGEX_SQLITE    = `^sqlite:\/\/(memory|\/?[\w.-]+(?:\/[\w.-]+)*(?:\.\w+)?)$` // sqlite:///path_to_db_file
	REGEX_MYSQL     = `^mysql:\/\/\w+:\w+@\w+(:\d+)?\/\w+$`                      // mysql://user:password@host:port/database
	REGEX_POSTGRES  = `^postgresql:\/\/\w+:\w+@\w+(:\d+)?\/\w+$`                 // postgresql://user:password@host:port/database
	REGEX_SQLSERVER = `sqlserver:\/\/\w+:\w+@\w+(:\d+)?\\[\w]+;database=\w+$`    // sqlserver://user:password@host:port\instance_name;database=database_name
)

type TYPE uint // The type of the database.

type Database struct {
	driver   gorm.Dialector
	provider *gorm.DB
}

func New(dsn string, dst ...interface{}) (*Database, error) {
	var driver gorm.Dialector

	parts := strings.Split(dsn, "://")
	if len(parts) != 2 {
		return nil, fmt.Errorf(ERROR_INVALID_DSN)
	}

	switch parts[0] {
	case TYPE_SQLITE:
		if matched, _ := regexp.MatchString(REGEX_SQLITE, dsn); !matched {
			return nil, fmt.Errorf(ERROR_INVALID_CONF, "sqlite", dsn)
		}
		if strings.Contains(dsn, "memory") {
			dsn = "file::memory:?cache=shared"
		}
		driver = sqlite.Open(parts[1])
	case TYPE_MYSQL:
		if matched, _ := regexp.MatchString(REGEX_MYSQL, dsn); !matched {
			return nil, fmt.Errorf(ERROR_INVALID_CONF, "mysql", dsn)
		}
		driver = mysql.Open(dsn)
	case TYPE_POSTGRES:
		if matched, _ := regexp.MatchString(REGEX_POSTGRES, dsn); !matched {
			return nil, fmt.Errorf(ERROR_INVALID_CONF, "postgres", dsn)
		}
		driver = postgres.Open(dsn)
	case TYPE_SQLSERVER:
		if matched, _ := regexp.MatchString(REGEX_SQLSERVER, dsn); !matched {
			return nil, fmt.Errorf(ERROR_INVALID_CONF, "sqlserver", dsn)
		}
		driver = sqlserver.Open(dsn)
	default:
		return nil, errors.New("Database type not supported please choose one of the following: sqlite, mysql, postgres, sqlserver")
	}

	provider, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if len(dst) > 0 {
		provider.AutoMigrate(dst...)
	}

	return &Database{driver, provider}, nil
}
