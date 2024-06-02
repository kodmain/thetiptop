package repositories_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestCreateClientRepository test de la création d'un client
func TestCreateClientRepository(t *testing.T) {
	// Création du mock SQL
	db, mock, err := sqlmock.New()
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

	// Création du repository client avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)

	// Données de transfert pour créer un client
	dto := &transfert.Client{
		Email:    "hello@world.com",
		Password: "password",
	}

	// Cas de création réussie
	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","validation_email","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"`).WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			nil,              // DeletedAt
			"hello@world.com",
			sqlmock.AnyArg(), // Password (hashed)
			false,            // ValidationEmail
			false,            // CGU
			false,            // Newsletter
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("b0d583fb-7d32-436f-9328-29620e8ca87b"))
		mock.ExpectCommit()

		entity, err := repo.Create(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas de contrainte unique échouée
	t.Run("unique constraint failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","validation_email","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"`).WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			nil,              // DeletedAt
			"hello@world.com",
			sqlmock.AnyArg(), // Password (hashed)
			false,            // ValidationEmail
			false,            // CGU
			false,            // Newsletter
		).WillReturnError(fmt.Errorf("UNIQUE constraint failed: clients.email"))
		mock.ExpectRollback()

		entity, err := repo.Create(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, errors.New("client already exists"), err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'autres erreurs lors de la création
	t.Run("other error during creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","validation_email","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"`).WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			nil,              // DeletedAt
			"hello@world.com",
			sqlmock.AnyArg(), // Password (hashed)
			false,            // ValidationEmail
			false,            // CGU
			false,            // Newsletter
		).WillReturnError(fmt.Errorf("some other error"))
		mock.ExpectRollback()

		entity, err := repo.Create(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "some other error", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// TestReadClientRepository test de la lecture d'un client
func TestReadClientRepository(t *testing.T) {
	// Création du mock SQL
	db, mock, err := sqlmock.New()
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

	// Création du repository client avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)

	// Données de transfert pour lire un client
	dto := &transfert.Client{
		Email: "hello@world.com",
	}

	// Cas de lecture réussie
	t.Run("successful read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."email" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs("hello@world.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("b0d583fb-7d32-436f-9328-29620e8ca87b", "hello@world.com"))

		entity, err := repo.Read(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, "hello@world.com", entity.Email)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où le client n'existe pas
	t.Run("client not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."email" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs("hello@world.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

		entity, err := repo.Read(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'autres erreurs lors de la lecture
	t.Run("other error during read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."email" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs("hello@world.com", 1).
			WillReturnError(errors.New("some other error"))

		entity, err := repo.Read(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "some other error", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
