package database

import "gorm.io/gorm"

// Option représente une fonction de configuration pour la requête GORM
type Option func(*gorm.DB) *gorm.DB

// GroupBy retourne une Option qui ajoute une clause GROUP BY
func GroupBy(column string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(column).Group(column)
	}
}

// Limit retourne une Option qui ajoute une clause LIMIT
func Limit(limit int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// Order retourne une Option qui ajoute une clause ORDER BY
func Order(order string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}
