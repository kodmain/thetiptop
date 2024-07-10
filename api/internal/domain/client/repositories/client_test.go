package repositories_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
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
		Email:    aws.String("hello@world.com"),
		Password: aws.String("password"),
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
		Email: aws.String("hello@world.com"),
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

func TestUpdateClientRepository(t *testing.T) {
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

	t.Run("successful update", func(t *testing.T) {
		// Ajout des attentes pour la transaction
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"cgu"=\$6,"newsletter"=\$7 WHERE "clients"."deleted_at" IS NULL AND "id" = \$8`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uuid).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateClient(&entities.Client{
			ID: uuid,
		})

		assert.Nil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestDeleteClientRepository(t *testing.T) {
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

	obj := &transfert.Client{
		Email: aws.String("hello@world.com"),
	}

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"."email" = \$2 AND "clients"."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), obj.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteClient(obj)

		assert.Nil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("delete with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"."email" = \$2 AND "clients"."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), "hello@world.com").
			WillReturnError(errors.New("some error"))
		mock.ExpectRollback()

		err := repo.DeleteClient(obj)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "some error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestCreateValidationRepository(t *testing.T) {
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

	// Création du repository validation avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)

	luhn := token.Generate(6)

	// Données de transfert pour créer une validation
	dto := &transfert.Validation{
		Token:    &luhn,
		ClientID: aws.String(uuid),
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "validations" \("id","created_at","updated_at","deleted_at","token","type","validated","client_id","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Token,        // Token
				sqlmock.AnyArg(), // Type
				false,            // Validated
				dto.ClientID,     // ClientID
				sqlmock.AnyArg(), // ExpiresAt
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		entity, err := repo.CreateValidation(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("creation with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "validations" \("id","created_at","updated_at","deleted_at","token","type","validated","client_id","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Token,        // Token
				sqlmock.AnyArg(), // Type
				false,            // Validated
				dto.ClientID,     // ClientID
				sqlmock.AnyArg(), // ExpiresAt
			).WillReturnError(errors.New("some error"))
		mock.ExpectRollback()

		entity, err := repo.CreateValidation(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "some error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestReadValidationRepository(t *testing.T) {
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

	// Création du repository validation avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)

	luhn := token.Generate(6)

	// Données de transfert pour lire une validation
	dto := &transfert.Validation{
		Token:    &luhn,
		ClientID: aws.String(uuid),
	}

	t.Run("successful read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "token", "client_id"}).AddRow("some-id", dto.Token, dto.ClientID))

		entity, err := repo.ReadValidation(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, dto.Token.String(), entity.Token.String())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("validation not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "token", "client_id"}))

		entity, err := repo.ReadValidation(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("other error during read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnError(errors.New("some other error"))

		entity, err := repo.ReadValidation(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "some other error", err.Error())

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUpdateValidationRepository(t *testing.T) {
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

	// Création du repository validation avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)
	luhn := token.Generate(6)

	entity := &entities.Validation{
		ID:        "some-id",
		Token:     &luhn,
		ClientID:  &uuid,
		Type:      entities.PasswordRecover,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"expires_at"=\$8 WHERE "validations"."deleted_at" IS NULL AND "id" = \$9`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, entity.ExpiresAt, entity.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateValidation(entity)

		assert.Nil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("update with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"expires_at"=\$8 WHERE "validations"."deleted_at" IS NULL AND "id" = \$9`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, entity.ExpiresAt, entity.ID).
			WillReturnError(errors.New("some error"))
		mock.ExpectRollback()

		err := repo.UpdateValidation(entity)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "some error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestDeleteValidationRepository(t *testing.T) {
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

	// Création du repository validation avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)

	luhn := token.Generate(6)
	dto := &transfert.Validation{
		Token:    &luhn,
		ClientID: aws.String("client-id"),
	}

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "deleted_at"=\$1 WHERE \("validations"\."token" = \$2 AND "validations"\."client_id" = \$3\) AND "validations"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Token, dto.ClientID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteValidation(dto)

		assert.Nil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("delete with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "deleted_at"=\$1 WHERE \("validations"\."token" = \$2 AND "validations"\."client_id" = \$3\) AND "validations"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Token, dto.ClientID).
			WillReturnError(errors.New("some error"))
		mock.ExpectRollback()

		err := repo.DeleteValidation(dto)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "some error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
