package repositories_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
)

func TestClientRepository(t *testing.T) {

	err := database.New(&database.Databases{
		"default": &database.Database{
			Protocol: database.SQLite,
			DBname:   ":memory:",
		},
	})

	assert.Nil(t, err)

	repo := repositories.NewClientRepository("default")
	assert.NotNil(t, repo)

	dto := &transfert.Client{
		Email:    "hello@world.com",
		Password: "password",
	}

	entity, err := repo.Create(dto)

	assert.Nil(t, err)
	assert.NotNil(t, entity)

	entity, err = repo.Create(dto)

	assert.NotNil(t, err)
	assert.Nil(t, entity)

	dto.Password = "" // Empty password because we cant test the BCRYPT hash

	entity, err = repo.Read(dto)

	assert.Nil(t, err)
	assert.NotNil(t, entity)
}
