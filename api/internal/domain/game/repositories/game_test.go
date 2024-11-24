package repositories_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setup() (*repositories.GameRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	dbInstance, _ := database.FromDB(gormDB)
	repo := repositories.NewGameRepository(dbInstance)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestCreateTicket(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Ticket{
		Prize: aws.String("PrizeA"),
		Token: aws.String("unique-token"),
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				dto.Token,        // Token
				dto.Prize,        // Prize
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		entity, err := repo.CreateTicket(dto)
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("creation with missing prize", func(t *testing.T) {
		dtoWithoutPrize := &transfert.Ticket{
			Token: aws.String("unique-token"),
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				dtoWithoutPrize.Token,
				nil, // Prize is missing
			).WillReturnError(fmt.Errorf("constraint violation"))

		mock.ExpectRollback()

		entity, err := repo.CreateTicket(dtoWithoutPrize)
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("creation with duplicate token", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				dto.Token,        // Token
				dto.Prize,        // Prize
			).WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))

		mock.ExpectRollback()

		entity, err := repo.CreateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("creation with database connection error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				dto.Token,        // Token
				dto.Prize,        // Prize
			).WillReturnError(fmt.Errorf("database is unavailable"))

		mock.ExpectRollback()

		entity, err := repo.CreateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful creation with custom options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				dto.Token,        // Token
				dto.Prize,        // Prize
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		entity, err := repo.CreateTicket(dto, database.Limit(1))
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateTickets(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	t.Run("successful creation of multiple tickets", func(t *testing.T) {
		tickets := []*transfert.Ticket{
			{Prize: aws.String("PrizeA"), Token: aws.String("TokenA")},
			{Prize: aws.String("PrizeB"), Token: aws.String("TokenB")},
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (Ticket 1)
				sqlmock.AnyArg(), // CreatedAt (Ticket 1)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 1)
				nil,              // DeletedAt (Ticket 1)
				nil,              // CredentialID (Ticket 1)
				"TokenA",         // Token (Ticket 1)
				"PrizeA",         // Prize (Ticket 1)

				sqlmock.AnyArg(), // ID (Ticket 2)
				sqlmock.AnyArg(), // CreatedAt (Ticket 2)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 2)
				nil,              // DeletedAt (Ticket 2)
				nil,              // CredentialID (Ticket 2)
				"TokenB",         // Token (Ticket 2)
				"PrizeB",         // Prize (Ticket 2)
			).WillReturnResult(sqlmock.NewResult(2, 2))
		mock.ExpectCommit()

		err := repo.CreateTickets(tickets)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("creation with duplicate token", func(t *testing.T) {
		tickets := []*transfert.Ticket{
			{Prize: aws.String("PrizeA"), Token: aws.String("TokenA")},
			{Prize: aws.String("PrizeB"), Token: aws.String("TokenB")},
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (Ticket 1)
				sqlmock.AnyArg(), // CreatedAt (Ticket 1)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 1)
				nil,              // DeletedAt (Ticket 1)
				nil,              // CredentialID (Ticket 1)
				"TokenA",         // Token (Ticket 1)
				"PrizeA",         // Prize (Ticket 1)

				sqlmock.AnyArg(), // ID (Ticket 2)
				sqlmock.AnyArg(), // CreatedAt (Ticket 2)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 2)
				nil,              // DeletedAt (Ticket 2)
				nil,              // CredentialID (Ticket 2)
				"TokenB",         // Token (Ticket 2)
				"PrizeB",         // Prize (Ticket 2)
			).WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))

		mock.ExpectRollback()

		err := repo.CreateTickets(tickets)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database unavailable", func(t *testing.T) {
		tickets := []*transfert.Ticket{
			{Prize: aws.String("PrizeA"), Token: aws.String("TokenA")},
			{Prize: aws.String("PrizeB"), Token: aws.String("TokenB")},
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (Ticket 1)
				sqlmock.AnyArg(), // CreatedAt (Ticket 1)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 1)
				nil,              // DeletedAt (Ticket 1)
				nil,              // CredentialID (Ticket 1)
				"TokenA",         // Token (Ticket 1)
				"PrizeA",         // Prize (Ticket 1)

				sqlmock.AnyArg(), // ID (Ticket 2)
				sqlmock.AnyArg(), // CreatedAt (Ticket 2)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 2)
				nil,              // DeletedAt (Ticket 2)
				nil,              // CredentialID (Ticket 2)
				"TokenB",         // Token (Ticket 2)
				"PrizeB",         // Prize (Ticket 2)
			).WillReturnError(fmt.Errorf("database is unavailable"))

		mock.ExpectRollback()

		err := repo.CreateTickets(tickets)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful creation with custom options", func(t *testing.T) {
		tickets := []*transfert.Ticket{
			{Prize: aws.String("PrizeA"), Token: aws.String("TokenA")},
			{Prize: aws.String("PrizeB"), Token: aws.String("TokenB")},
		}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "tickets" \("id","created_at","updated_at","deleted_at","credential_id","token","prize"\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (Ticket 1)
				sqlmock.AnyArg(), // CreatedAt (Ticket 1)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 1)
				nil,              // DeletedAt (Ticket 1)
				nil,              // CredentialID (Ticket 1)
				"TokenA",         // Token (Ticket 1)
				"PrizeA",         // Prize (Ticket 1)

				sqlmock.AnyArg(), // ID (Ticket 2)
				sqlmock.AnyArg(), // CreatedAt (Ticket 2)
				sqlmock.AnyArg(), // UpdatedAt (Ticket 2)
				nil,              // DeletedAt (Ticket 2)
				nil,              // CredentialID (Ticket 2)
				"TokenB",         // Token (Ticket 2)
				"PrizeB",         // Prize (Ticket 2),
			).WillReturnResult(sqlmock.NewResult(2, 2))
		mock.ExpectCommit()

		err := repo.CreateTickets(tickets, database.Limit(1))
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReadTicket(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Ticket{
		Prize: aws.String("PrizeA"),
		Token: aws.String("unique-token"),
	}

	t.Run("successful read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE \("tickets"\."prize" = \$1 AND "tickets"\."token" = \$2\) AND "tickets"\."deleted_at" IS NULL ORDER BY "tickets"\."id" LIMIT \$3`).
			WithArgs(dto.Prize, dto.Token, 1). // Inclure la limite
			WillReturnRows(sqlmock.NewRows([]string{"id", "prize", "token"}).AddRow("some-id", "PrizeA", "unique-token"))

		entity, err := repo.ReadTicket(dto)

		// Vérification que la requête ne retourne pas d'erreur
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		// Vérification des champs retournés
		if assert.NotNil(t, entity.Prize) { // Vérifie avant d'accéder à Prize
			assert.Equal(t, *dto.Prize, *entity.Prize)
		}
		if assert.NotNil(t, entity.Token) { // Vérifie avant d'accéder à Token
			assert.Equal(t, "unique-token", entity.Token.String())
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE \("tickets"\."prize" = \$1 AND "tickets"\."token" = \$2\) AND "tickets"\."deleted_at" IS NULL ORDER BY "tickets"\."id" LIMIT \$3`).
			WithArgs(dto.Prize, dto.Token, 1).          // Inclure la limite
			WillReturnRows(sqlmock.NewRows([]string{})) // Aucune ligne retournée

		entity, err := repo.ReadTicket(dto)

		// Vérification qu'une erreur est retournée
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "ticket.not_found", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE \("tickets"\."prize" = \$1 AND "tickets"\."token" = \$2\) AND "tickets"\."deleted_at" IS NULL ORDER BY "tickets"\."id" LIMIT \$3`).
			WithArgs(dto.Prize, dto.Token, 1).
			WillReturnError(fmt.Errorf("database error"))

		entity, err := repo.ReadTicket(dto, database.Limit(1))

		// Vérification de l'erreur et de la nullité de l'entité retournée
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReadTickets(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Ticket{
		Prize: aws.String("PrizeA"),
		Token: aws.String("unique-token"),
	}

	t.Run("successful read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE`).
			WithArgs(dto.Prize, dto.Token).
			WillReturnRows(sqlmock.NewRows([]string{"id", "prize", "token"}).
				AddRow("ticket1", "PrizeA", "unique-token").
				AddRow("ticket2", "PrizeA", "unique-token"))

		tickets, err := repo.ReadTickets(dto)
		assert.Nil(t, err)
		assert.NotNil(t, tickets)
		assert.Len(t, tickets, 2)
		assert.Equal(t, "PrizeA", *tickets[0].Prize)
		assert.Equal(t, "unique-token", tickets[0].Token.String())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no tickets found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE`).
			WithArgs(dto.Prize, dto.Token).
			WillReturnRows(sqlmock.NewRows([]string{}))

		tickets, err := repo.ReadTickets(dto)
		assert.Nil(t, err)
		assert.NotNil(t, tickets)
		assert.Len(t, tickets, 0)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE`).
			WithArgs(dto.Prize, dto.Token).
			WillReturnError(fmt.Errorf("database error"))

		tickets, err := repo.ReadTickets(dto)
		assert.NotNil(t, err)
		assert.Nil(t, tickets)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successful read with custom options", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE`).
			WithArgs(dto.Prize, dto.Token).
			WillReturnRows(sqlmock.NewRows([]string{"id", "prize", "token"}).
				AddRow("ticket1", "PrizeA", "unique-token").
				AddRow("ticket2", "PrizeA", "unique-token"))

		tickets, err := repo.ReadTickets(dto, database.Order("created_at DESC"))
		assert.Nil(t, err)
		assert.NotNil(t, tickets)
		assert.Len(t, tickets, 2)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateTicket(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	token1 := token.Generate(12)

	entity := &entities.Ticket{
		ID:           "some-id",
		Prize:        aws.String("PrizeA"),
		Token:        token1,
		CredentialID: nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("successful update with options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tickets" SET`).
			WithArgs(
				sqlmock.AnyArg(),    // CreatedAt
				sqlmock.AnyArg(),    // UpdatedAt
				nil,                 // DeletedAt
				entity.CredentialID, // CredentialID
				entity.Token,        // Token
				entity.Prize,        // Prize
				entity.ID,           // ID
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateTicket(entity, database.Order("created_at DESC"))
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update failure with options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tickets" SET`).
			WithArgs(
				sqlmock.AnyArg(),    // CreatedAt
				sqlmock.AnyArg(),    // UpdatedAt
				nil,                 // DeletedAt
				entity.CredentialID, // CredentialID
				entity.Token,        // Token
				entity.Prize,        // Prize
				entity.ID,           // ID
			).WillReturnError(fmt.Errorf("update error"))
		mock.ExpectRollback()

		err := repo.UpdateTicket(entity, database.Order("created_at DESC"))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteTicket(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Ticket{
		Token: aws.String("delete-token"),
	}

	t.Run("successful deletion with options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"=\$1 WHERE "tickets"."token" = \$2 AND "tickets"."deleted_at" IS NULL`).
			WithArgs(
				sqlmock.AnyArg(), // Timestamp pour "deleted_at"
				dto.Token,        // Token
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteTicket(dto, database.Limit(1))
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("deletion failure with options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"=\$1 WHERE "tickets"."token" = \$2 AND "tickets"."deleted_at" IS NULL`).
			WithArgs(
				sqlmock.AnyArg(), // Timestamp pour "deleted_at"
				dto.Token,        // Token
			).WillReturnError(fmt.Errorf("deletion error"))
		mock.ExpectRollback()

		err := repo.DeleteTicket(dto, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCountTicket(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Ticket{
		Prize: aws.String("PrizeA"),
	}

	t.Run("successful count with options", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE`).
			WithArgs(dto.Prize).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

		count, err := repo.CountTicket(dto, database.Order("ASC"))
		assert.Nil(t, err)
		assert.Equal(t, 42, count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("count failure with options", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE`).
			WithArgs(dto.Prize).
			WillReturnError(fmt.Errorf("count error"))

		count, err := repo.CountTicket(dto, database.Order("ASC"))
		assert.NotNil(t, err)
		assert.Equal(t, 0, count)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
