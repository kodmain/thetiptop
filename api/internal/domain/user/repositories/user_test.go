package repositories_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
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

// TestCreateUserRepository test de la création d'un client
func TestCreateUserRepository(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

	// Création du mock SQL
	db, _, err := sqlmock.New()
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

	// Création du repository user avec l'instance de base de données mockée
	repo := repositories.NewUserRepository(dbInstance)
	assert.NotNil(t, repo)
}

func setup() (*repositories.UserRepository, sqlmock.Sqlmock, *sql.DB) {
	config.Load(aws.String("../../../../config.test.yml"))

	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	dbInstance, _ := database.FromDB(gormDB)
	repo := repositories.NewUserRepository(dbInstance)

	return repo, mock, db
}

func TestCreateCredential(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	dto := &transfert.Credential{
		Email:    aws.String("hello@world.com"),
		Password: aws.String("password"),
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction de l'instruction SQL pour supprimer la colonne client_id qui n'existe pas dans la requête réelle
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				dto.Email,
				sqlmock.AnyArg(), // Password (hashed)
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		entity, err := repo.CreateCredential(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unique constraint failed", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction de l'instruction SQL pour supprimer la colonne client_id
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				dto.Email,
				sqlmock.AnyArg(),
			).WillReturnError(fmt.Errorf("UNIQUE constraint failed: credentials.email"))

		mock.ExpectRollback()

		entity, err := repo.CreateCredential(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("random error", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction de l'instruction SQL pour supprimer la colonne client_id
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				dto.Email,
				sqlmock.AnyArg(),
			).WillReturnError(fmt.Errorf("random-error"))

		mock.ExpectRollback()

		entity, err := repo.CreateCredential(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestReadCredential test de la lecture des credentials basé sur l'email
func TestReadCredential(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour lire un credential par email
	dto := &transfert.Credential{
		Email: aws.String("hello@world.com"),
	}

	// Cas de lecture réussie
	t.Run("successful read", func(t *testing.T) {
		// Correction de l'expression régulière pour inclure la clause ORDER BY
		mock.ExpectQuery(`SELECT \* FROM "credentials" WHERE "credentials"\."email" = \$1 AND "credentials"\."deleted_at" IS NULL ORDER BY "credentials"\."id" LIMIT \$2`).
			WithArgs(dto.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow("some-uuid", dto.Email, "hashed-password"))

		// Appel de la méthode ReadCredential du repository
		entity, err := repo.ReadCredential(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, dto.Email, entity.Email)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où le credential n'existe pas
	t.Run("credential not found", func(t *testing.T) {
		dto := &transfert.Credential{
			Email: aws.String("non-existing-email@world.com"),
		}

		mock.ExpectQuery(`SELECT \* FROM "credentials" WHERE "credentials"\."email" = \$1 AND "credentials"\."deleted_at" IS NULL ORDER BY "credentials"\."id" LIMIT \$2`).
			WithArgs(dto.Email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}))

		entity, err := repo.ReadCredential(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "credential.not_found")

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où le credential n'existe pas
	t.Run("credential not found", func(t *testing.T) {
		dto := &transfert.Credential{
			Email: aws.String("non-existing-email@world.com"),
		}

		mock.ExpectQuery(`SELECT \* FROM "credentials" WHERE "credentials"\."email" = \$1 AND "credentials"\."deleted_at" IS NULL ORDER BY "credentials"\."id" LIMIT \$2`).
			WithArgs(dto.Email, 1).
			WillReturnError(fmt.Errorf("some client error"))

		entity, err := repo.ReadCredential(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "common.internal_error")

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestUpdateCredential teste la mise à jour des credentials
func TestUpdateCredential(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	entity := &entities.Credential{
		ID:       uuid,
		Email:    aws.String("updated@world.com"),
		Password: aws.String("new-hashed-password"),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction de l'instruction SQL : suppression de la colonne `client_id`
		mock.ExpectExec(`UPDATE "credentials" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5 WHERE "credentials"\."deleted_at" IS NULL AND "id" = \$6`).
			WithArgs(
				sqlmock.AnyArg(), // created_at (générée automatiquement)
				sqlmock.AnyArg(), // updated_at (générée automatiquement)
				nil,              // deleted_at (NULL)
				entity.Email,     // mise à jour de l'email
				entity.Password,  // mise à jour du mot de passe
				entity.ID,        // ID du credential
			).WillReturnResult(sqlmock.NewResult(1, 1)) // Résultat de succès (1 ligne affectée)

		mock.ExpectCommit()

		err := repo.UpdateCredential(entity)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update failure", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction de l'instruction SQL : suppression de la colonne `client_id`
		mock.ExpectExec(`UPDATE "credentials" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5 WHERE "credentials"\."deleted_at" IS NULL AND "id" = \$6`).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				nil,              // deleted_at
				entity.Email,     // mise à jour de l'email
				entity.Password,  // mise à jour du mot de passe
				entity.ID,        // ID du credential
			).WillReturnError(fmt.Errorf("some update error"))

		mock.ExpectRollback()

		err := repo.UpdateCredential(entity)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestDeleteCredential teste la suppression logique des credentials (soft delete)
func TestDeleteCredential(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour supprimer un credential
	dto := &transfert.Credential{
		Email: aws.String("delete@world.com"),
	}

	// Cas de suppression réussie (soft delete)
	t.Run("successful deletion", func(t *testing.T) {
		// Mock de la requête pour supprimer un credential (soft delete)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "credentials" SET "deleted_at"=\$1 WHERE "credentials"\."email" = \$2 AND "credentials"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Email).    // La date actuelle sera utilisée pour "deleted_at"
			WillReturnResult(sqlmock.NewResult(1, 1)) // 1 ligne affectée par la suppression
		mock.ExpectCommit()

		// Appel de la méthode DeleteCredential du repository
		err := repo.DeleteCredential(dto)

		// Vérification des résultats
		assert.Nil(t, err)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où la suppression échoue
	t.Run("deletion failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "credentials" SET "deleted_at"=\$1 WHERE "credentials"\."email" = \$2 AND "credentials"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Email).
			WillReturnError(fmt.Errorf("some delete error"))
		mock.ExpectRollback()

		// Appel de la méthode DeleteCredential du repository
		err := repo.DeleteCredential(dto)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCreateValidation teste la création d'une validation
func TestCreateValidation(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	luhn := token.Generate(6)
	dto := &transfert.Validation{
		Token:    luhn.PointerString(),
		ClientID: aws.String("client-uuid"),
	}

	// Cas de création réussie
	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "validations" \("id","created_at","updated_at","deleted_at","token","type","validated","client_id","employee_id","credential_id","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Token,        // Token
				sqlmock.AnyArg(), // Type
				false,            // Validated
				dto.ClientID,     // ClientID
				nil,              // EmployeeID (probablement NULL)
				nil,              // CredentialID (probablement NULL)
				sqlmock.AnyArg(), // ExpiresAt
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		entity, err := repo.CreateValidation(dto)

		assert.Nil(t, err)
		assert.NotNil(t, entity)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où la création échoue
	t.Run("creation with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "validations" \("id","created_at","updated_at","deleted_at","token","type","validated","client_id","employee_id","credential_id","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Token,        // Token
				sqlmock.AnyArg(), // Type
				false,            // Validated
				dto.ClientID,     // ClientID
				nil,              // EmployeeID (probablement NULL)
				nil,              // CredentialID (probablement NULL)
				sqlmock.AnyArg(), // ExpiresAt
			).WillReturnError(fmt.Errorf("some error"))
		mock.ExpectRollback()

		entity, err := repo.CreateValidation(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "common.internal_error")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

// TestReadValidation teste la lecture d'une validation
func TestReadValidation(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Génération d'un token de validation
	luhn := token.Generate(6)
	dto := &transfert.Validation{
		Token:    luhn.PointerString(),
		ClientID: aws.String("client-uuid"),
	}

	// Données simulées pour une entité Validation
	entity := &entities.Validation{
		ID:        "some-id",
		Token:     &luhn,
		ClientID:  aws.String("client-uuid"),
		Type:      entities.PasswordRecover,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Cas de lecture réussie
	t.Run("successful read", func(t *testing.T) {
		// Mock de la requête SQL avec les bons arguments, y compris la limite
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "token", "client_id"}).AddRow(entity.ID, *dto.Token, *dto.ClientID))

		// Appel de la méthode ReadValidation du repository
		result, err := repo.ReadValidation(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, *dto.Token, result.Token.String())
		assert.Equal(t, *dto.ClientID, *result.ClientID)

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où la validation n'est pas trouvée
	t.Run("validation not found", func(t *testing.T) {
		// Mock pour simuler le cas où aucune ligne n'est retournée
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnRows(sqlmock.NewRows([]string{})) // Pas de résultat

		// Appel de la méthode ReadValidation du repository
		result, err := repo.ReadValidation(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "validation.not_found")

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas d'une erreur inattendue lors de la lecture
	t.Run("read with error", func(t *testing.T) {
		// Mock pour simuler une erreur SQL
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnError(fmt.Errorf("some error")) // Simuler une erreur

		// Appel de la méthode ReadValidation du repository
		result, err := repo.ReadValidation(dto)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "common.internal_error")

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestUpdateValidation teste la méthode UpdateValidation du UserRepository
func TestUpdateValidation(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Génération d'un token et d'une entité de validation
	luhn := token.Generate(6)
	entity := &entities.Validation{
		ID:        "some-id",
		Token:     &luhn,
		ClientID:  aws.String("client-uuid"),
		Type:      entities.PasswordRecover,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Cas de mise à jour réussie
	t.Run("successful update", func(t *testing.T) {
		// Mock de la requête SQL pour la mise à jour de l'entité
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"employee_id"=\$8,"credential_id"=\$9,"expires_at"=\$10 WHERE "validations"."deleted_at" IS NULL AND "id" = \$11`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, nil, nil, entity.ExpiresAt, entity.ID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Succès de la mise à jour
		mock.ExpectCommit()

		// Appel de la méthode UpdateValidation
		err := repo.UpdateValidation(entity)

		// Vérification des résultats
		assert.Nil(t, err) // Pas d'erreur attendue

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où l'entité est nil
	t.Run("nil entity", func(t *testing.T) {
		// Appel de la méthode avec une entité nil
		err := repo.UpdateValidation(nil)

		// Vérification que l'erreur est correcte
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		// Pas d'interaction SQL attendue, donc aucune attente SQL à vérifier
	})

	// Cas d'échec lors de la mise à jour
	t.Run("update failure", func(t *testing.T) {
		// Mock pour simuler une erreur SQL lors de la mise à jour
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"employee_id"=\$8,"credential_id"=\$9,"expires_at"=\$10 WHERE "validations"."deleted_at" IS NULL AND "id" = \$11`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, nil, nil, entity.ExpiresAt, entity.ID).
			WillReturnError(fmt.Errorf("update failed")) // Simuler une erreur
		mock.ExpectRollback()

		// Appel de la méthode UpdateValidation
		err := repo.UpdateValidation(entity)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.EqualError(t, err, "common.internal_error")

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
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
	repo := repositories.NewUserRepository(dbInstance)

	// Génération d'un token et d'une entité Validation à supprimer
	luhn := token.Generate(6)
	dto := &transfert.Validation{
		Token:    luhn.PointerString(),
		ClientID: aws.String("client-id"),
	}

	t.Run("successful delete", func(t *testing.T) {
		// Mock de la requête SQL pour l'UPDATE (soft delete)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "deleted_at"=\$1 WHERE \("validations"\."token" = \$2 AND "validations"\."client_id" = \$3\) AND "validations"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Token, dto.ClientID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Succès de la suppression
		mock.ExpectCommit()

		// Appel de la méthode DeleteValidation
		err := repo.DeleteValidation(dto)

		// Vérification des résultats
		assert.Nil(t, err) // Pas d'erreur attendue

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete with error", func(t *testing.T) {
		// Mock pour simuler une erreur SQL lors de l'UPDATE
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "deleted_at"=\$1 WHERE \("validations"\."token" = \$2 AND "validations"\."client_id" = \$3\) AND "validations"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Token, dto.ClientID).
			WillReturnError(fmt.Errorf("some error")) // Simuler une erreur
		mock.ExpectRollback()

		// Appel de la méthode DeleteValidation
		err := repo.DeleteValidation(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.EqualError(t, err, "common.internal_error")

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateClient(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour créer un client
	dto := &transfert.Client{
		CGU: aws.Bool(true),
	}

	// Cas de création réussie
	t.Run("successful creation", func(t *testing.T) {
		// Démarrage de la transaction
		mock.ExpectBegin()

		// Insertion dans la table clients avec la colonne credential_id ajoutée
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","credential_id","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				true,             // CGU
				false,            // Newsletter
			).WillReturnResult(sqlmock.NewResult(1, 1))

		// Validation de la transaction
		mock.ExpectCommit()

		// Appeler la fonction CreateClient à tester
		entity, err := repo.CreateClient(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas d'autres erreurs lors de la création
	t.Run("other error during creation", func(t *testing.T) {
		mock.ExpectBegin()

		// Corriger l'expression régulière pour inclure credential_id
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","credential_id","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (UUID)
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				nil,              // CredentialID
				true,             // CGU
				false,            // Newsletter
			).WillReturnError(fmt.Errorf("some other error"))

		mock.ExpectRollback()

		// Appel à la méthode testée
		entity, err := repo.CreateClient(dto)

		// Assertions pour vérifier le comportement attendu
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		// Vérification des expectations SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestReadClient test de la lecture d'un client basé sur les attributs Client
func TestReadClient(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour lire un client
	dto := &transfert.Client{
		ID: aws.String(uuid),
	}

	// Cas de lecture réussie
	t.Run("successful read", func(t *testing.T) {
		// Mock de la requête pour lire un client par ID
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.ID, 1). // Assurez-vous que les arguments correspondent à ceux utilisés dans votre méthode
			WillReturnRows(sqlmock.NewRows([]string{"id", "cgu", "newsletter"}).AddRow(uuid, true, false))

		// Appel de la méthode ReadClient du repository
		entity, err := repo.ReadClient(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, *dto.ID, entity.ID)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où le client n'existe pas
	t.Run("client not found", func(t *testing.T) {
		dto := &transfert.Client{
			ID: aws.String(uuid),
		}

		// Mise à jour du mock SQL pour inclure la clause ORDER BY
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{})) // Pas de résultat retourné

		entity, err := repo.ReadClient(dto)

		// Vérification de l'erreur : le client n'existe pas, donc gorm.ErrRecordNotFound est attendu
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "client.not_found")

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestReadValidations(t *testing.T) {
	// Fonction setup pour initialiser les dépendances
	repo, mock, db := setup()
	defer db.Close()

	// Cas de lecture réussie
	t.Run("successful read", func(t *testing.T) {
		validationType := "password_recovery"
		dto := &transfert.Validation{
			Type: &validationType,
		}

		// Données simulées pour les validations
		expectedValidations := []*entities.Validation{
			{
				ID:        "validation-id-1",
				Type:      entities.PasswordRecover,
				ClientID:  aws.String("client-id-1"),
				Validated: false,
			},
			{
				ID:        "validation-id-2",
				Type:      entities.PasswordRecover,
				ClientID:  aws.String("client-id-2"),
				Validated: true,
			},
		}

		// Mock de la requête SQL
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE "validations"\."type" = \$1 AND "validations"\."deleted_at" IS NULL`).
			WithArgs(validationType).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type", "client_id", "validated"}).
				AddRow(expectedValidations[0].ID, expectedValidations[0].Type, *expectedValidations[0].ClientID, expectedValidations[0].Validated).
				AddRow(expectedValidations[1].ID, expectedValidations[1].Type, *expectedValidations[1].ClientID, expectedValidations[1].Validated))

		// Appel de la méthode
		result, err := repo.ReadValidations(dto)

		// Vérifications
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedValidations[0].ID, result[0].ID)
		assert.Equal(t, expectedValidations[1].ID, result[1].ID)

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où aucune validation n'est trouvée
	t.Run("no validations found", func(t *testing.T) {
		validationType := "email_verification"
		dto := &transfert.Validation{
			Type: &validationType,
		}

		// Mock pour simuler aucune ligne trouvée
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE "validations"\."type" = \$1 AND "validations"\."deleted_at" IS NULL`).
			WithArgs(validationType).
			WillReturnRows(sqlmock.NewRows([]string{}))

		// Appel de la méthode
		result, err := repo.ReadValidations(dto)

		// Vérifications
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0) // Résultat vide attendu

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où une erreur SQL survient
	t.Run("database error", func(t *testing.T) {
		validationType := "email_verification"
		dto := &transfert.Validation{
			Type: &validationType,
		}

		// Mock pour simuler une erreur SQL
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE "validations"\."type" = \$1 AND "validations"\."deleted_at" IS NULL`).
			WithArgs(validationType).
			WillReturnError(fmt.Errorf("database error"))

		// Appel de la méthode
		result, err := repo.ReadValidations(dto)

		// Vérifications
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "common.internal_error")

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestUpdateClient teste la mise à jour des clients
func TestUpdateClient(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	entity := &entities.Client{
		ID:         "b0d583fb-7d32-436f-9328-29620e8ca87b",
		Newsletter: aws.Bool(true),
		CGU:        aws.Bool(true),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction : ajout de la colonne `credential_id` dans l'instruction SQL
		mock.ExpectExec(`UPDATE "clients" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"credential_id"=\$4,"cgu"=\$5,"newsletter"=\$6 WHERE "clients"\."deleted_at" IS NULL AND "id" = \$7`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at
				sqlmock.AnyArg(),  // updated_at
				nil,               // deleted_at
				nil,               // credential_id
				entity.CGU,        // mise à jour de CGU
				entity.Newsletter, // mise à jour de la newsletter
				entity.ID,         // ID du client
			).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Succès (1 ligne affectée)

		mock.ExpectCommit()

		err := repo.UpdateClient(entity)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update failure", func(t *testing.T) {
		mock.ExpectBegin()

		// Correction : ajout de la colonne `credential_id`
		mock.ExpectExec(`UPDATE "clients" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"credential_id"=\$4,"cgu"=\$5,"newsletter"=\$6 WHERE "clients"\."deleted_at" IS NULL AND "id" = \$7`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at
				sqlmock.AnyArg(),  // updated_at
				nil,               // deleted_at
				nil,               // credential_id
				entity.CGU,        // mise à jour de CGU
				entity.Newsletter, // mise à jour de la newsletter
				entity.ID,         // ID du client
			).WillReturnError(fmt.Errorf("some update error"))

		mock.ExpectRollback()

		err := repo.UpdateClient(entity)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestDeleteClient teste la suppression logique des clients (soft delete)
func TestDeleteClient(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour supprimer un client
	dto := &transfert.Client{
		ID: aws.String("client-uuid"),
	}

	// Cas de suppression réussie (soft delete)
	t.Run("successful deletion", func(t *testing.T) {
		// Mock de la requête pour supprimer un client (soft delete)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"\."id" = \$2 AND "clients"\."id" = \$3 AND "clients"\."deleted_at" IS NULL`).
			WithArgs(
				sqlmock.AnyArg(), // Date de suppression
				dto.ID,           // ID du client
				dto.ID,           // Deuxième condition sur l'ID
			).WillReturnResult(sqlmock.NewResult(1, 1)) // 1 ligne affectée par la suppression
		mock.ExpectCommit()

		// Appel de la méthode DeleteClient du repository
		err := repo.DeleteClient(dto)

		// Vérification des résultats
		assert.Nil(t, err)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Cas où la suppression échoue
	t.Run("deletion failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"\."id" = \$2 AND "clients"\."id" = \$3 AND "clients"\."deleted_at" IS NULL`).
			WithArgs(
				sqlmock.AnyArg(), // Date de suppression
				dto.ID,           // Premier argument ID
				dto.ID,           // Deuxième argument ID
			).
			WillReturnError(fmt.Errorf("some delete error"))
		mock.ExpectRollback()

		// Appel de la méthode DeleteClient du repository
		err := repo.DeleteClient(dto)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateEmployee(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour créer un employé
	dto := &transfert.Employee{
		CredentialID: aws.String("credential-uuid"),
	}

	t.Run("successful creation", func(t *testing.T) {
		mock.ExpectBegin()

		// Insertion dans la table employees avec la colonne credential_id
		mock.ExpectExec(`INSERT INTO "employees" \("id","created_at","updated_at","deleted_at","credential_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5\)`).
			WithArgs(
				sqlmock.AnyArg(),  // ID (UUID)
				sqlmock.AnyArg(),  // CreatedAt
				sqlmock.AnyArg(),  // UpdatedAt
				nil,               // DeletedAt
				"credential-uuid", // CredentialID (mis à jour pour refléter la valeur correcte)
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Appel de la fonction CreateEmployee à tester
		entity, err := repo.CreateEmployee(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error during creation", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`INSERT INTO "employees" \("id","created_at","updated_at","deleted_at","credential_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5\)`).
			WithArgs(
				sqlmock.AnyArg(),  // ID (UUID)
				sqlmock.AnyArg(),  // CreatedAt
				sqlmock.AnyArg(),  // UpdatedAt
				nil,               // DeletedAt
				"credential-uuid", // CredentialID (mis à jour pour refléter la valeur correcte)
			).WillReturnError(fmt.Errorf("creation error"))

		mock.ExpectRollback()

		entity, err := repo.CreateEmployee(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "common.internal_error", err.Error())

		// Vérification des attentes
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReadEmployee(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	dto := &transfert.Employee{
		ID: aws.String(uuid),
	}

	t.Run("successful read", func(t *testing.T) {
		// Ajustement de l'expression régulière pour matcher la requête SQL exacte
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE "employees"\."id" = \$1 AND "employees"\."deleted_at" IS NULL AND "employees"\."id" = \$2 ORDER BY "employees"\."id" LIMIT \$3`).
			WithArgs(dto.ID, dto.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "credential_id"}).AddRow(uuid, "credential-uuid"))

		entity, err := repo.ReadEmployee(dto)

		// Vérification que la requête n'a pas retourné d'erreur
		assert.Nil(t, err)

		// Vérification que l'entité n'est pas nulle avant de faire des assertions supplémentaires
		if assert.NotNil(t, entity) {
			assert.Equal(t, *dto.ID, entity.ID)
			assert.Equal(t, "credential-uuid", *entity.CredentialID)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("employee not found", func(t *testing.T) {
		// Même ajustement pour la requête "not found"
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE "employees"\."id" = \$1 AND "employees"\."deleted_at" IS NULL AND "employees"\."id" = \$2 ORDER BY "employees"\."id" LIMIT \$3`).
			WithArgs(dto.ID, dto.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "credential_id"}))

		entity, err := repo.ReadEmployee(dto)

		// Si l'employé n'est pas trouvé, une erreur doit être renvoyée
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.EqualError(t, err, "employee.not_found")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateEmployee(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	entity := &entities.Employee{
		ID:           "b0d583fb-7d32-436f-9328-29620e8ca87b",
		CredentialID: aws.String("credential-uuid"),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "employees" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"credential_id"=\$4 WHERE "employees"\."deleted_at" IS NULL AND "id" = \$5`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at
				sqlmock.AnyArg(),  // updated_at
				nil,               // deleted_at
				"credential-uuid", // CredentialID
				entity.ID,         // ID de l'employé
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.UpdateEmployee(entity)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("update failure", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(`UPDATE "employees" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"credential_id"=\$4 WHERE "employees"\."deleted_at" IS NULL AND "id" = \$5`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at
				sqlmock.AnyArg(),  // updated_at
				nil,               // deleted_at
				"credential-uuid", // CredentialID
				entity.ID,         // ID de l'employé
			).WillReturnError(fmt.Errorf("update error"))

		mock.ExpectRollback()

		err := repo.UpdateEmployee(entity)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteEmployee(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	dto := &transfert.Employee{
		ID: aws.String("employee-uuid"),
	}

	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "employees" SET "deleted_at"=\$1 WHERE "employees"\."id" = \$2 AND "employees"\."id" = \$3 AND "employees"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.ID, dto.ID). // Deux fois l'ID correspondant à deux conditions
			WillReturnResult(sqlmock.NewResult(1, 1))   // 1 ligne affectée par la suppression
		mock.ExpectCommit()

		err := repo.DeleteEmployee(dto)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("deletion failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "employees" SET "deleted_at"=\$1 WHERE "employees"\."id" = \$2 AND "employees"\."id" = \$3 AND "employees"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.ID, dto.ID).
			WillReturnError(fmt.Errorf("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteEmployee(dto)
		assert.NotNil(t, err)
		assert.Equal(t, "common.internal_error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReadUser(t *testing.T) {

	t.Run("user is a client", func(t *testing.T) {
		repo, mock, db := setup()
		defer db.Close()

		// Données de transfert pour lire un utilisateur
		dto := &transfert.User{
			ID: aws.String("user-uuid"),
		}

		// Ajuster la requête pour matcher le comportement actuel
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.ToClient().ID, 1). // Le deuxième argument est la limite
			WillReturnRows(sqlmock.NewRows([]string{"id", "credential_id"}).AddRow("client-uuid", "credential-uuid"))

		client, employee, err := repo.ReadUser(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, client)
		assert.Nil(t, employee)
		assert.Equal(t, "client-uuid", client.ID)

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user is an employee", func(t *testing.T) {
		repo, mock, db := setup()
		defer db.Close()

		// Données de transfert pour lire un utilisateur
		dto := &transfert.User{
			ID: aws.String("user-uuid"),
		}

		// Mock pour simuler qu'aucun client n'est trouvé
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.ToClient().ID, 1).             // Le deuxième argument est la limite
			WillReturnRows(sqlmock.NewRows([]string{})) // Aucun client trouvé

		// Mock pour simuler qu'un employé est trouvé
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE "employees"\."id" = \$1 AND "employees"\."deleted_at" IS NULL ORDER BY "employees"\."id" LIMIT \$2`).
			WithArgs(dto.ToEmployee().ID, 1). // Le deuxième argument est la limite
			WillReturnRows(sqlmock.NewRows([]string{"id", "credential_id"}).AddRow("employee-uuid", "credential-uuid"))

		// Mock de la requête pour chercher les validations associées à l'employé
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE employee_id = \$1 AND "validations"\."deleted_at" IS NULL`).
			WithArgs("employee-uuid").
			WillReturnRows(sqlmock.NewRows([]string{})) // Simuler aucune validation trouvée

		// Appel de la méthode ReadUser
		client, employee, err := repo.ReadUser(dto)

		// Vérification des résultats
		assert.Nil(t, err)         // Aucune erreur attendue
		assert.NotNil(t, employee) // Employee ne doit pas être nil
		assert.Nil(t, client)      // Client doit être nil, car il s'agit d'un employé

		// Vérification des champs de l'employé
		assert.Equal(t, "employee-uuid", employee.ID)

		// Vérification si CredentialID n'est pas nil avant d'accéder à la valeur
		if employee.CredentialID != nil {
			assert.Equal(t, "credential-uuid", *employee.CredentialID)
		}

		// Vérification des attentes SQL
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		repo, mock, db := setup()
		defer db.Close()

		dto := &transfert.User{
			ID: aws.String("user-uuid"),
		}

		// Mock pour simuler qu'aucun client n'est trouvé
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL LIMIT \$2`).
			WithArgs(dto.ToClient().ID, 1).             // Le deuxième argument est la limite dynamique
			WillReturnRows(sqlmock.NewRows([]string{})) // Pas de client trouvé

		// Mock pour simuler qu'aucun employé n'est trouvé
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE "employees"\."id" = \$1 AND "employees"\."deleted_at" IS NULL LIMIT \$2`).
			WithArgs(dto.ToEmployee().ID, 1).           // Le deuxième argument est la limite dynamique
			WillReturnRows(sqlmock.NewRows([]string{})) // Pas d'employé trouvé

		client, employee, err := repo.ReadUser(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.Nil(t, client)
		assert.Nil(t, employee)
		assert.EqualError(t, err, "user.not_found")
	})

	t.Run("error on reading client and employee", func(t *testing.T) {
		repo, mock, db := setup()
		defer db.Close()

		dto := &transfert.User{
			ID: aws.String("user-uuid"),
		}

		// Simuler une erreur lors de la lecture des clients
		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL LIMIT 1`).
			WithArgs(dto.ToClient().ID).
			WillReturnError(fmt.Errorf("some client error"))

		// Simuler une erreur lors de la lecture des employés
		mock.ExpectQuery(`SELECT \* FROM "employees" WHERE "employees"\."id" = \$1 AND "employees"\."deleted_at" IS NULL LIMIT 1`).
			WithArgs(dto.ToClient().ID).
			WillReturnError(fmt.Errorf("some employee error"))

		client, employee, err := repo.ReadUser(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.Nil(t, client)
		assert.Nil(t, employee)
		assert.EqualError(t, err, "user.not_found")
	})
}
