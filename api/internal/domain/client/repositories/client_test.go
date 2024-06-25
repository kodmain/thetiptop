package repositories_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	uuid string = "b0d583fb-7d32-436f-9328-29620e8ca87b"
)

// TestCreateClientRepository test de la création d'un client
func TestCreateClientRepository(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

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
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				false,            // CGU
				false,            // Newsletter
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`INSERT INTO "validations" \("id","created_at","updated_at","deleted_at","token","type","validated","client_id","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				sqlmock.AnyArg(), // Token
				sqlmock.AnyArg(), // Type
				false,            // Validated
				sqlmock.AnyArg(), // ClientID
				sqlmock.AnyArg(), // ExpiresAt
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Appeler la fonction CreateClient à tester
		entity, err := repo.CreateClient(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas de contrainte unique échouée
	t.Run("unique constraint failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				false,            // CGU
				false,            // Newsletter
			).WillReturnError(fmt.Errorf("UNIQUE constraint failed: clients.email"))
		mock.ExpectRollback()

		entity, err := repo.CreateClient(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, errors.New("client already exists"), err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'autres erreurs lors de la création
	t.Run("other error during creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","email","password","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				false,            // CGU
				false,            // Newsletter
			).WillReturnError(fmt.Errorf("some other error"))
		mock.ExpectRollback()

		entity, err := repo.CreateClient(dto)

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
			WithArgs(dto.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(uuid, dto.Email))
		entity, err := repo.ReadClient(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, dto.Email, entity.Email)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où le client n'existe pas
	t.Run("client not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."email" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

		entity, err := repo.ReadClient(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'autres erreurs lors de la lecture
	t.Run("other error during read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."email" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.Email, 1).
			WillReturnError(errors.New("some other error"))

		entity, err := repo.ReadClient(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "some other error", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
