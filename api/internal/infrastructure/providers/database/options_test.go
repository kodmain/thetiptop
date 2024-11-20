package database

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID    uint
	Name  string
	Age   int
	Group string
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Création d'une table de test
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		return nil, err
	}

	// Insertion de données
	db.Create(&TestModel{ID: 1, Name: "Alice", Age: 30, Group: "A"})
	db.Create(&TestModel{ID: 2, Name: "Bob", Age: 25, Group: "B"})
	db.Create(&TestModel{ID: 3, Name: "Charlie", Age: 35, Group: "A"})
	db.Create(&TestModel{ID: 4, Name: "David", Age: 40, Group: "C"})

	return db, nil
}

func TestGroupBy(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	var results []struct {
		Group string
		Count int
	}

	// Utiliser GroupBy pour regrouper les enregistrements par "group"
	query := GroupBy("`group`")(db.Table("test_models").Select("`group`, COUNT(*) as count")).Scan(&results)

	if query.Error != nil {
		t.Fatalf("Failed to execute GroupBy: %v", query.Error)
	}

	expectedGroups := 3 // Group A, B, C
	if len(results) != expectedGroups {
		t.Errorf("Expected %d groups, got %d", expectedGroups, len(results))
	}

	// Vérifier les noms des groupes (optionnel)
	expectedGroupNames := map[string]bool{
		"A": true,
		"B": true,
		"C": true,
	}

	for _, result := range results {
		if _, exists := expectedGroupNames[result.Group]; !exists {
			t.Errorf("Unexpected group name: %s", result.Group)
		}
	}
}

func TestWhere(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	var results []TestModel
	query := Where("age > ?", 30)(db).Find(&results)

	if query.Error != nil {
		t.Fatalf("Failed to execute Where: %v", query.Error)
	}

	expectedCount := 2 // Charlie and David
	if len(results) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(results))
	}
}

func TestLimit(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	var results []TestModel
	query := Limit(2)(db).Find(&results)

	if query.Error != nil {
		t.Fatalf("Failed to execute Limit: %v", query.Error)
	}

	expectedCount := 2
	if len(results) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(results))
	}
}

func TestOrder(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	var results []TestModel
	query := Order("age DESC")(db).Find(&results)

	if query.Error != nil {
		t.Fatalf("Failed to execute Order: %v", query.Error)
	}

	if len(results) > 0 && results[0].Age != 40 { // David has the highest age
		t.Errorf("Expected first result to be David (age 40), got age %d", results[0].Age)
	}
}
