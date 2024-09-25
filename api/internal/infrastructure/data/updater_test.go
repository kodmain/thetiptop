package data_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

type Address struct {
	Street string
	City   string
	Zip    string
}

type Company struct {
	Name    string
	Website string
}

type ComplexEntity struct {
	ID             string
	Name           *string // Pointeur pour tester entityField.Kind() == reflect.Ptr
	Age            int     // Non pointeur pour tester reflect.DeepEqual
	IsActive       *bool   // Pointeur pour tester entityField.Kind() == reflect.Ptr
	Emails         []string
	Address        Address
	Company        *Company
	PhoneNumbers   map[string]string
	FavoriteColors []string
}

type ComplexEntityDTO struct {
	ID             string
	Name           *string
	Age            *int
	IsActive       *bool
	Emails         *[]string
	Address        *Address
	Company        *Company
	PhoneNumbers   *map[string]string
	FavoriteColors *[]string
}

func TestUpdateComplexEntityWithDto(t *testing.T) {

	// Test pour vérifier la mise à jour des pointeurs (entityField.Kind() == reflect.Ptr)
	t.Run("successful update with pointer assignment", func(t *testing.T) {
		// Initialisation de l'entité ComplexEntity avec des champs pointeurs
		originalName := "John Doe"
		originalIsActive := true
		entity := &ComplexEntity{
			ID:       "123",
			Name:     &originalName,     // Champ pointeur pour tester entityField.Kind() == reflect.Ptr
			IsActive: &originalIsActive, // Champ pointeur pour tester entityField.Kind() == reflect.Ptr
		}

		// Initialisation du DTO ComplexEntityDTO avec un champ pointeur
		newName := "Jane Doe"
		newIsActive := false
		dto := &ComplexEntityDTO{
			ID:       "123",
			Name:     &newName,     // Pointeur pour tester l'assignation directe du pointeur
			IsActive: &newIsActive, // Pointeur pour tester l'assignation directe du pointeur
		}

		// Appel de la méthode UpdateEntityWithDto
		data.UpdateEntityWithDto(entity, dto)

		// Vérification que les pointeurs sont mis à jour
		assert.Equal(t, "Jane Doe", *entity.Name) // Le pointeur du champ Name est mis à jour
		assert.Equal(t, false, *entity.IsActive)  // Le pointeur du champ IsActive est mis à jour
	})

	// Test pour vérifier la mise à jour lorsque les champs non pointeurs sont différents (reflect.DeepEqual)
	t.Run("successful update when fields are different", func(t *testing.T) {
		// Initialisation de l'entité ComplexEntity
		entity := &ComplexEntity{
			ID:     "123",
			Age:    30, // Non pointeur pour tester reflect.DeepEqual
			Emails: []string{"john@example.com"},
			Address: Address{
				Street: "123 Main St",
				City:   "Metropolis",
				Zip:    "12345",
			},
		}

		// Initialisation du DTO ComplexEntityDTO avec des valeurs modifiées
		newAge := 35 // Nouveau champ pour tester la mise à jour via reflect.DeepEqual
		newEmails := []string{"jane@example.com"}
		newAddress := Address{
			Street: "456 Main St",
			City:   "Smallville",
			Zip:    "54321",
		}
		dto := &ComplexEntityDTO{
			ID:      "123",
			Age:     &newAge,     // Pour tester la différence entre les valeurs de l'entité et du DTO
			Emails:  &newEmails,  // Liste mise à jour pour forcer reflect.DeepEqual
			Address: &newAddress, // Adresse modifiée
		}

		// Appel de la méthode UpdateEntityWithDto
		data.UpdateEntityWithDto(entity, dto)

		// Vérification que reflect.DeepEqual force la mise à jour des champs
		assert.Equal(t, 35, entity.Age)                              // Age mis à jour
		assert.Equal(t, []string{"jane@example.com"}, entity.Emails) // Emails mis à jour
		assert.Equal(t, Address{
			Street: "456 Main St",
			City:   "Smallville",
			Zip:    "54321",
		}, entity.Address) // Address mis à jour
	})

	// Test pour garantir que rien n'est mis à jour si les champs sont identiques (reflect.DeepEqual)
	t.Run("no update when fields are equal", func(t *testing.T) {
		// Initialisation de l'entité ComplexEntity
		originalName := "John Doe"
		entity := &ComplexEntity{
			ID:   "123",
			Name: &originalName,
		}

		// DTO avec les mêmes valeurs que l'entité
		sameName := "John Doe"
		dto := &ComplexEntityDTO{
			ID:   "123",
			Name: &sameName, // Même nom pour tester la non-mise à jour
		}

		// Appel de la méthode UpdateEntityWithDto
		data.UpdateEntityWithDto(entity, dto)

		// Vérification : Il ne doit pas y avoir de changement
		assert.Equal(t, "John Doe", *entity.Name) // Name inchangé
	})

	// Nouveau test pour forcer l'entrée dans reflect.DeepEqual
	t.Run("update when non-pointer fields are different", func(t *testing.T) {
		// Initialisation de l'entité ComplexEntity avec des valeurs de base
		originalName := "John Doe"
		entity := &ComplexEntity{
			ID:       "123",
			Name:     &originalName, // Champ pointeur
			Age:      30,            // Champ non pointeur
			IsActive: new(bool),     // Champ pointeur
			Emails:   []string{"john@example.com"},
			Address: Address{
				Street: "123 Main St",
				City:   "Metropolis",
				Zip:    "12345",
			},
		}

		// Initialisation du DTO avec des valeurs différentes
		newName := "Jane Doe"                     // Différent de l'entité
		newAge := 35                              // Différent de l'entité
		newIsActive := true                       // Différent de l'entité
		newEmails := []string{"jane@example.com"} // Différent de l'entité
		newAddress := Address{
			Street: "456 Main St",
			City:   "Smallville",
			Zip:    "54321",
		}
		dto := &ComplexEntityDTO{
			ID:       "123",        // ID identique
			Name:     &newName,     // Différent pour forcer l'entrée dans reflect.DeepEqual
			Age:      &newAge,      // Différent pour forcer la mise à jour
			IsActive: &newIsActive, // Différent pour forcer la mise à jour
			Emails:   &newEmails,   // Différent pour forcer la mise à jour
			Address:  &newAddress,  // Différent pour forcer la mise à jour
		}

		// Appel de la méthode UpdateEntityWithDto pour tester reflect.DeepEqual
		data.UpdateEntityWithDto(entity, dto)

		// Vérification des résultats : les champs doivent être mis à jour car ils sont différents
		assert.Equal(t, "Jane Doe", *entity.Name)                    // Name mis à jour
		assert.Equal(t, 35, entity.Age)                              // Age mis à jour
		assert.Equal(t, true, *entity.IsActive)                      // IsActive mis à jour
		assert.Equal(t, []string{"jane@example.com"}, entity.Emails) // Emails mis à jour
		assert.Equal(t, Address{
			Street: "456 Main St",
			City:   "Smallville",
			Zip:    "54321",
		}, entity.Address) // Address mis à jour
	})

	// Test pour entrer dans le bloc où entityField n'est pas un pointeur mais dtoField l'est
	t.Run("successful update when entity field is not pointer but dto field is", func(t *testing.T) {
		// Initialisation de l'entité ComplexEntity avec des valeurs non-pointeur
		entity := &ComplexEntity{
			ID:     "123",
			Age:    40, // Champ non pointeur
			Emails: []string{"john@example.com"},
		}

		// Initialisation du DTO avec des valeurs modifiées et des pointeurs
		newAge := 50
		newEmails := []string{"jane@example.com"}
		dto := &ComplexEntityDTO{
			ID:     "123",
			Age:    &newAge,    // Champ pointeur dans le DTO
			Emails: &newEmails, // Champ pointeur dans le DTO
		}

		// Appel de la méthode UpdateEntityWithDto pour tester l'assignation après dé-référencement
		data.UpdateEntityWithDto(entity, dto)

		// Vérification que les champs sont mis à jour correctement
		assert.Equal(t, 50, entity.Age)                              // Age mis à jour
		assert.Equal(t, []string{"jane@example.com"}, entity.Emails) // Emails mis à jour
	})
}

type SimpleEntity struct {
	ID     string
	Age    int     // Non pointeur pour tester reflect.DeepEqual
	Salary float64 // Non pointeur pour tester reflect.DeepEqual
}

type SimpleEntityDTO struct {
	ID     string
	Age    int     // Non pointeur
	Salary float64 // Non pointeur
}

func TestUpdateEntityWithDto(t *testing.T) {
	t.Run("successful update when fields are different", func(t *testing.T) {
		// Initialisation de l'entité SimpleEntity
		entity := &SimpleEntity{
			ID:     "123",
			Age:    25,      // Valeur actuelle de l'entité
			Salary: 50000.0, // Valeur actuelle
		}

		// Initialisation du DTO SimpleEntityDTO avec des valeurs modifiées
		dto := &SimpleEntityDTO{
			ID:     "123",   // ID non modifié
			Age:    30,      // Age différent pour forcer l'entrée dans reflect.DeepEqual
			Salary: 60000.0, // Salaire différent pour tester reflect.DeepEqual
		}

		// Appel de la méthode UpdateEntityWithDto
		data.UpdateEntityWithDto(entity, dto)

		// Vérification que les champs non pointeurs sont mis à jour via reflect.DeepEqual
		assert.Equal(t, 30, entity.Age)         // Age mis à jour
		assert.Equal(t, 60000.0, entity.Salary) // Salaire mis à jour
	})
}
