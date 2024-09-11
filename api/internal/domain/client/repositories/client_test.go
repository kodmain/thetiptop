package repositories_test

import (
	"database/sql"
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

	// Création du repository client avec l'instance de base de données mockée
	repo := repositories.NewClientRepository(dbInstance)
	assert.NotNil(t, repo)
}

func setup() (*repositories.ClientRepository, sqlmock.Sqlmock, *sql.DB) {
	config.Load(aws.String("../../../../config.test.yml"))

	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	dbInstance, _ := database.FromDB(gormDB)
	repo := repositories.NewClientRepository(dbInstance)

	return repo, mock, db
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

		// Insertion dans la table clients
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
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
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'autres erreurs lors de la création
	t.Run("other error during creation", func(t *testing.T) {
		mock.ExpectBegin()
		// Corriger l'expression régulière pour correspondre aux placeholders PostgreSQL ($1, $2, etc.)
		mock.ExpectExec(`INSERT INTO "clients" \("id","created_at","updated_at","deleted_at","cgu","newsletter"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID (UUID)
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				true,             // CGU
				false,            // Newsletter
			).WillReturnError(fmt.Errorf("some other error"))
		mock.ExpectRollback()

		// Appel à la méthode testée
		entity, err := repo.CreateClient(dto)

		// Assertions pour vérifier le comportement attendu
		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, "some other error", err.Error())

		// Vérification des expectations SQL
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
			WithArgs(dto.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cgu", "newsletter"}).AddRow(uuid, true, false))

		// Appel de la méthode ReadClient du repository
		entity, err := repo.ReadClient(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, *dto.ID, entity.ID)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où le client n'existe pas
	t.Run("client not found", func(t *testing.T) {
		dto := &transfert.Client{
			ID: aws.String(uuid),
		}

		mock.ExpectQuery(`SELECT \* FROM "clients" WHERE "clients"\."id" = \$1 AND "clients"\."deleted_at" IS NULL ORDER BY "clients"\."id" LIMIT \$2`).
			WithArgs(dto.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cgu", "newsletter"}))

		entity, err := repo.ReadClient(dto)

		assert.NotNil(t, err)
		assert.Nil(t, entity)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestCreateCredential(t *testing.T) {
	repo, mock, db := setup()
	defer db.Close()

	// Données de transfert pour créer un client
	dto := &transfert.Credential{
		Email:    aws.String("hello@world.com"),
		Password: aws.String("password"),
	}

	// Cas de création réussie
	t.Run("successful creation", func(t *testing.T) {
		// Démarrage de la transaction
		mock.ExpectBegin()

		// Insertion dans la table credentials (ajout de client_id)
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password","client_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				nil,              // ClientID (parce qu'il est probablement nul)
			).WillReturnResult(sqlmock.NewResult(1, 1))

		// Validation de la transaction
		mock.ExpectCommit()

		// Appeler la fonction CreateCredential à tester
		entity, err := repo.CreateCredential(dto)

		// Vérification des résultats
		assert.Nil(t, err)
		assert.NotNil(t, entity)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("unique constraint failed", func(t *testing.T) {
		// Démarrage de la transaction
		mock.ExpectBegin()

		// Insertion dans la table credentials
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password","client_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				nil,              // ClientID (parce qu'il est probablement nul)
			).WillReturnError(fmt.Errorf("UNIQUE constraint failed: credentials.email")) // Simuler l'erreur de contrainte UNIQUE

		// Annulation de la transaction
		mock.ExpectRollback()

		// Appeler la fonction CreateCredential à tester
		entity, err := repo.CreateCredential(dto)

		// Vérification des résultats
		assert.NotNil(t, err) // On s'attend à une erreur
		assert.Nil(t, entity) // L'entité ne doit pas être créée
		assert.Equal(t, "credential already exists", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("unique constraint failed", func(t *testing.T) {
		// Démarrage de la transaction
		mock.ExpectBegin()

		// Insertion dans la table credentials
		mock.ExpectExec(`INSERT INTO "credentials" \("id","created_at","updated_at","deleted_at","email","password","client_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)`).
			WithArgs(
				sqlmock.AnyArg(), // ID
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				nil,              // DeletedAt
				dto.Email,        // Email
				sqlmock.AnyArg(), // Password (hashed)
				nil,              // ClientID (parce qu'il est probablement nul)
			).WillReturnError(fmt.Errorf("random-error")) // random error

		// Annulation de la transaction
		mock.ExpectRollback()

		// Appeler la fonction CreateCredential à tester
		entity, err := repo.CreateCredential(dto)

		// Vérification des résultats
		assert.NotNil(t, err) // On s'attend à une erreur
		assert.Nil(t, entity) // L'entité ne doit pas être créée
		assert.Equal(t, "random-error", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// TestUpdateClient teste la mise à jour des clients
func TestUpdateClient(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de l'entité Client à mettre à jour
	entity := &entities.Client{
		ID:         "b0d583fb-7d32-436f-9328-29620e8ca87b",
		Newsletter: aws.Bool(true),
		CGU:        aws.Bool(true),
	}

	// Cas de mise à jour réussie
	t.Run("successful update", func(t *testing.T) {
		// Mock de la requête pour la mise à jour du client
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"cgu"=\$4,"newsletter"=\$5 WHERE "clients"\."deleted_at" IS NULL AND "id" = \$6`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at (valeur temporelle générée automatiquement)
				sqlmock.AnyArg(),  // updated_at (valeur temporelle générée automatiquement)
				nil,               // deleted_at (car NULL)
				entity.CGU,        // mise à jour de CGU
				entity.Newsletter, // mise à jour de la newsletter
				entity.ID,         // ID du client
			).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Résultat de succès (1 ligne affectée)
		mock.ExpectCommit()

		// Appel de la méthode UpdateClient du repository
		err := repo.UpdateClient(entity)

		// Vérification des résultats
		assert.Nil(t, err)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où la mise à jour échoue
	t.Run("update failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"cgu"=\$4,"newsletter"=\$5 WHERE "clients"\."deleted_at" IS NULL AND "id" = \$6`).
			WithArgs(
				sqlmock.AnyArg(),  // created_at
				sqlmock.AnyArg(),  // updated_at
				nil,               // deleted_at
				entity.CGU,        // mise à jour de CGU
				entity.Newsletter, // mise à jour de la newsletter
				entity.ID,         // ID du client
			).
			WillReturnError(fmt.Errorf("some update error"))
		mock.ExpectRollback()

		// Appel de la méthode UpdateClient du repository
		err := repo.UpdateClient(entity)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Equal(t, "some update error", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"\."id" = \$2 AND "clients"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.ID).       // La date actuelle sera utilisée pour "deleted_at"
			WillReturnResult(sqlmock.NewResult(1, 1)) // 1 ligne affectée par la suppression
		mock.ExpectCommit()

		// Appel de la méthode DeleteClient du repository
		err := repo.DeleteClient(dto)

		// Vérification des résultats
		assert.Nil(t, err)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où la suppression échoue
	t.Run("deletion failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "clients" SET "deleted_at"=\$1 WHERE "clients"\."id" = \$2 AND "clients"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.ID).
			WillReturnError(fmt.Errorf("some delete error"))
		mock.ExpectRollback()

		// Appel de la méthode DeleteClient du repository
		err := repo.DeleteClient(dto)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Equal(t, "some delete error", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// TestUpdateCredential teste la mise à jour des credentials
func TestUpdateCredential(t *testing.T) {
	// Initialisation du repository, du mock et de la base de données
	repo, mock, db := setup()
	defer db.Close()

	// Données de l'entité Credential à mettre à jour
	entity := &entities.Credential{
		ID:       uuid,
		Email:    aws.String("updated@world.com"),
		Password: aws.String("new-hashed-password"),
	}

	// Cas de mise à jour réussie
	t.Run("successful update", func(t *testing.T) {
		// Mock de la requête pour la mise à jour du credential
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "credentials" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"client_id"=\$6 WHERE "credentials"\."deleted_at" IS NULL AND "id" = \$7`).
			WithArgs(
				sqlmock.AnyArg(), // created_at (probablement une valeur temporelle générée automatiquement)
				sqlmock.AnyArg(), // updated_at (valeur temporelle générée automatiquement)
				nil,              // deleted_at (car NULL)
				entity.Email,     // mise à jour de l'email
				entity.Password,  // mise à jour du mot de passe
				nil,              // client_id (car probablement NULL)
				entity.ID,        // ID du credential
			).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Résultat de succès (1 ligne affectée)
		mock.ExpectCommit()

		// Appel de la méthode UpdateCredential du repository
		err := repo.UpdateCredential(entity)

		// Vérification des résultats
		assert.Nil(t, err)

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où la mise à jour échoue
	t.Run("update failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "credentials" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"email"=\$4,"password"=\$5,"client_id"=\$6 WHERE "credentials"\."deleted_at" IS NULL AND "id" = \$7`).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				nil,              // deleted_at
				entity.Email,     // mise à jour de l'email
				entity.Password,  // mise à jour du mot de passe
				nil,              // client_id
				entity.ID,        // ID du credential
			).
			WillReturnError(fmt.Errorf("some update error"))
		mock.ExpectRollback()

		// Appel de la méthode UpdateCredential du repository
		err := repo.UpdateCredential(entity)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Equal(t, "some update error", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		assert.Equal(t, "some delete error", err.Error())

		// Vérification des attentes
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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

	// Cas où la création échoue
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
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
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
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// Vérification des attentes SQL
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas d'une erreur inattendue lors de la lecture
	t.Run("read with error", func(t *testing.T) {
		// Mock pour simuler une erreur SQL
		mock.ExpectQuery(`SELECT \* FROM "validations" WHERE \("validations"\."token" = \$1 AND "validations"\."client_id" = \$2\) AND "validations"\."deleted_at" IS NULL ORDER BY "validations"\."id" LIMIT \$3`).
			WithArgs(dto.Token, dto.ClientID, 1).
			WillReturnError(errors.New("some error")) // Simuler une erreur

		// Appel de la méthode ReadValidation du repository
		result, err := repo.ReadValidation(dto)

		// Vérification que l'erreur est bien renvoyée
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "some error")

		// Vérification des attentes SQL
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// TestUpdateValidation teste la méthode UpdateValidation du ClientRepository
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
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"expires_at"=\$8 WHERE "validations"."deleted_at" IS NULL AND "id" = \$9`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, entity.ExpiresAt, entity.ID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Succès de la mise à jour
		mock.ExpectCommit()

		// Appel de la méthode UpdateValidation
		err := repo.UpdateValidation(entity)

		// Vérification des résultats
		assert.Nil(t, err) // Pas d'erreur attendue

		// Vérification des attentes SQL
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Cas où l'entité est nil
	t.Run("nil entity", func(t *testing.T) {
		// Appel de la méthode avec une entité nil
		err := repo.UpdateValidation(nil)

		// Vérification que l'erreur est correcte
		assert.NotNil(t, err)
		assert.Equal(t, "invalid value, should be pointer to struct or slice", err.Error())

		// Pas d'interaction SQL attendue, donc aucune attente SQL à vérifier
	})

	// Cas d'échec lors de la mise à jour
	t.Run("update failure", func(t *testing.T) {
		// Mock pour simuler une erreur SQL lors de la mise à jour
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"token"=\$4,"type"=\$5,"validated"=\$6,"client_id"=\$7,"expires_at"=\$8 WHERE "validations"."deleted_at" IS NULL AND "id" = \$9`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, entity.Token, sqlmock.AnyArg(), sqlmock.AnyArg(), entity.ClientID, entity.ExpiresAt, entity.ID).
			WillReturnError(errors.New("update failed")) // Simuler une erreur
		mock.ExpectRollback()

		// Appel de la méthode UpdateValidation
		err := repo.UpdateValidation(entity)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.EqualError(t, err, "update failed")

		// Vérification des attentes SQL
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
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("delete with error", func(t *testing.T) {
		// Mock pour simuler une erreur SQL lors de l'UPDATE
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "validations" SET "deleted_at"=\$1 WHERE \("validations"\."token" = \$2 AND "validations"\."client_id" = \$3\) AND "validations"\."deleted_at" IS NULL`).
			WithArgs(sqlmock.AnyArg(), dto.Token, dto.ClientID).
			WillReturnError(errors.New("some error")) // Simuler une erreur
		mock.ExpectRollback()

		// Appel de la méthode DeleteValidation
		err := repo.DeleteValidation(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.EqualError(t, err, "some error")

		// Vérification des attentes SQL
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
