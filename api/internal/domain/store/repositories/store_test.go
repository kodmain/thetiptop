package repositories_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	errors_domain_store "github.com/kodmain/thetiptop/api/internal/domain/store/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setup initializes the test environment, creates a sqlmock database and returns a StoreRepository instance, the mock, and a cleanup function
// Parameters:
// - None
//
// Returns:
// - *repositories.StoreRepository: the repository instance
// - sqlmock.Sqlmock: the mock instance
// - func(): a cleanup function to close the DB
func setup() (*repositories.StoreRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	dbInstance, _ := database.FromDB(gormDB)
	repo := repositories.NewStoreRepository(dbInstance)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

// Test_CreateStores tests the CreateStores method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_CreateStores(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	stores := []*transfert.Store{
		{ID: aws.String("store-1")},
		{ID: aws.String("store-2")},
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		// GORM will insert (id, created_at, updated_at, deleted_at, label, is_online)
		mock.ExpectExec(`INSERT INTO "stores"`).
			WithArgs(
				sqlmock.AnyArg(), // ID for store-1
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // Label (nil)
				nil,              // IsOnline
				sqlmock.AnyArg(), // ID for store-2
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // Label (nil)
				nil,              // IsOnline
			).WillReturnResult(sqlmock.NewResult(2, 2))
		mock.ExpectCommit()

		err := repo.CreateStores(stores)
		assert.Nil(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "stores"`).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.CreateStores(stores, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_ReadStores tests the ReadStores method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_ReadStores(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Store{}

	t.Run("successful read", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE "stores"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label", "is_online"}).
				AddRow("s1", "store-1", nil).
				AddRow("s2", "store-2", true),
			)

		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."store_id" = \$1 AND "caisses"\."deleted_at" IS NULL`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).
				AddRow("c1", "s1").
				AddRow("c2", "s1"))

		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."store_id" = \$1 AND "caisses"\."deleted_at" IS NULL`).
			WithArgs("s2").
			WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}))

		stores, err := repo.ReadStores(dto)
		assert.Nil(t, err)
		assert.Len(t, stores, 2)
		assert.NotNil(t, stores[0].Label)
		assert.Equal(t, "store-1", *stores[0].Label)
		assert.Len(t, stores[0].Caisses, 2)
		assert.Equal(t, "c1", stores[0].Caisses[0].ID)
		assert.Equal(t, "c2", stores[0].Caisses[1].ID)

		assert.NotNil(t, stores[1].Label)
		assert.Equal(t, "store-2", *stores[1].Label)
		assert.Len(t, stores[1].Caisses, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no data found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE "stores"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{}))

		stores, err := repo.ReadStores(dto)
		assert.Nil(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE "stores"\."deleted_at" IS NULL`).
			WillReturnError(errors.New("db error"))

		stores, err := repo.ReadStores(dto, database.Limit(1))
		assert.NotNil(t, err)
		assert.Nil(t, stores)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_ReadStore tests the ReadStore method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_ReadStore(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Store{ID: aws.String("store-123")}

	t.Run("store found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE ((?:"stores"\."id" = \$1 AND "stores"\."deleted_at" IS NULL)|("stores"\."deleted_at" IS NULL AND "stores"\."id" = \$1")) ORDER BY "stores"\."id" LIMIT \$2`).
			WithArgs("store-123", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "label", "is_online"}).
				AddRow("store-123", "my-store", nil),
			)

		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."store_id" = \$1 AND "caisses"\."deleted_at" IS NULL`).
			WithArgs("store-123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}))

		store, err := repo.ReadStore(dto)
		assert.Nil(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, "my-store", *store.Label)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("store not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE ((?:"stores"\."id" = \$1 AND "stores"\."deleted_at" IS NULL)|("stores"\."deleted_at" IS NULL AND "stores"\."id" = \$1")) ORDER BY "stores"\."id" LIMIT \$2`).
			WithArgs("store-456", 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		dtoNotFound := &transfert.Store{ID: aws.String("store-456")}
		store, err := repo.ReadStore(dtoNotFound)
		assert.Nil(t, store)
		assert.NotNil(t, err)
		assert.Equal(t, errors_domain_store.ErrStoreNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "stores" WHERE ((?:"stores"\."id" = \$1 AND "stores"\."deleted_at" IS NULL)|("stores"\."deleted_at" IS NULL AND "stores"\."id" = \$1")) ORDER BY "stores"\."id" LIMIT \$2`).
			WithArgs("store-123", 1).
			WillReturnError(errors.New("db error"))

		store, err := repo.ReadStore(dto, database.Limit(1))
		assert.Nil(t, store)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_DeleteStores tests the DeleteStores method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_DeleteStores(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	objs := []*transfert.Store{
		{ID: aws.String("store-1")},
		{ID: aws.String("store-2")},
	}

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "stores" SET "deleted_at"=`).
			WithArgs(sqlmock.AnyArg(), "store-1", "store-2").
			WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectCommit()

		err := repo.DeleteStores(objs)
		assert.Nil(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "stores"`).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.DeleteStores(objs, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_UpdateStores tests the UpdateStores method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_UpdateStores(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	lbl1 := "updated-store-1"
	lbl2 := "updated-store-2"
	objs := []*entities.Store{
		{ID: "store-1", Label: &lbl1},
		{ID: "store-2", Label: &lbl2},
	}

	t.Run("successful update", func(t *testing.T) {
		// Pour store-1
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "stores" SET .* WHERE "stores"\."deleted_at" IS NULL AND "id" = \$3`).
			WithArgs(sqlmock.AnyArg(), "updated-store-1", "store-1").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Pour store-2
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "stores" SET .* WHERE "stores"\."deleted_at" IS NULL AND "id" = \$3`).
			WithArgs(sqlmock.AnyArg(), "updated-store-2", "store-2").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateStores(objs)
		assert.Nil(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "stores"`).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.UpdateStores(objs, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_CreateCaisse tests the CreateCaisse method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_CreateCaisse(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Caisse{}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "caisses"`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				nil,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		caisse, err := repo.CreateCaisse(dto)
		assert.Nil(t, err)
		assert.NotNil(t, caisse)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "caisses"`).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		caisse, err := repo.CreateCaisse(dto, database.Limit(1))
		assert.NotNil(t, err)
		assert.Nil(t, caisse)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_ReadCaisse tests the ReadCaisse method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_ReadCaisse(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Caisse{ID: aws.String("caisse-123")}

	t.Run("caisse found", func(t *testing.T) {
		// La requête utilise deux arguments : l'ID et la limite
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."id" = \$1 AND "caisses"\."deleted_at" IS NULL ORDER BY "caisses"\."id" LIMIT \$2`).
			WithArgs("caisse-123", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).
				AddRow("caisse-123", nil))

		caisse, err := repo.ReadCaisse(dto)
		assert.Nil(t, err)
		assert.NotNil(t, caisse)
		assert.Equal(t, "caisse-123", caisse.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no caisse found", func(t *testing.T) {
		dtoNotFound := &transfert.Caisse{ID: aws.String("not-exist")}
		// De même, on s'attend à deux arguments : l'ID et la limite
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."id" = \$1 AND "caisses"\."deleted_at" IS NULL ORDER BY "caisses"\."id" LIMIT \$2`).
			WithArgs("not-exist", 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		caisse, err := repo.ReadCaisse(dtoNotFound)
		assert.NotNil(t, err)
		// Comme aucune ligne n'est trouvée, la fonction renvoie nil, nil (selon la logique actuelle)
		assert.Nil(t, caisse)
		assert.Equal(t, "caisse.not_found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."id" = \$1 AND "caisses"\."deleted_at" IS NULL ORDER BY "caisses"\."id" LIMIT \$2`).
			WithArgs("caisse-123", 1).
			WillReturnError(errors.New("db error"))

		caisse, err := repo.ReadCaisse(dto, database.Limit(1))
		assert.Nil(t, caisse)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_ReadCaisses tests the ReadCaisses method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_ReadCaisses(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Caisse{}

	t.Run("caisses found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "store_id"}).
				AddRow("c1", nil).
				AddRow("c2", nil))

		caisses, err := repo.ReadCaisses(dto)
		assert.Nil(t, err)
		assert.Len(t, caisses, 2)
		assert.Equal(t, "c1", caisses[0].ID)
		assert.Equal(t, "c2", caisses[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no caisses found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."deleted_at" IS NULL`).
			WillReturnRows(sqlmock.NewRows([]string{}))

		caisses, err := repo.ReadCaisses(dto)
		assert.Nil(t, err)
		assert.NotNil(t, caisses)
		assert.Len(t, caisses, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "caisses" WHERE "caisses"\."deleted_at" IS NULL`).
			WillReturnError(errors.New("db error"))

		caisses, err := repo.ReadCaisses(dto, database.Limit(1))
		assert.NotNil(t, err)
		assert.Nil(t, caisses)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_DeleteCaisse tests the DeleteCaisse method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_DeleteCaisse(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	dto := &transfert.Caisse{ID: aws.String("caisse-123")}

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectBegin()
		// On matche la requête exactement telle qu'elle est générée par GORM
		mock.ExpectExec(`UPDATE "caisses" SET "deleted_at"=\$1 WHERE "caisses"\."id" = \$2 AND "caisses"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), "caisse-123").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteCaisse(dto)
		assert.Nil(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "caisses" SET "deleted_at"=\$1 WHERE "caisses"\."id" = \$2 AND "caisses"\."deleted_at" IS NULL`).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.DeleteCaisse(dto, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Test_UpdateCaisse tests the UpdateCaisse method of StoreRepository
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_UpdateCaisse(t *testing.T) {
	repo, mock, cleanup := setup()
	defer cleanup()

	obj := &entities.Caisse{ID: "caisse-123", StoreID: aws.String("store-456")}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()
		// On utilise une regex large pour ne pas imposer l'ordre exact des champs dans le SET.
		// On sait que GORM met updated_at d'abord, store_id ensuite, donc on adapte WithArgs.
		mock.ExpectExec(`UPDATE "caisses" SET .* WHERE "caisses"\."deleted_at" IS NULL AND "id" = \$3`).
			WithArgs(sqlmock.AnyArg(), "store-456", "caisse-123"). // $1=updated_at(anyArg), $2=store-456, $3=caisse-123
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateCaisse(obj)
		assert.Nil(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "caisses"`).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.UpdateCaisse(obj, database.Limit(1))
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
