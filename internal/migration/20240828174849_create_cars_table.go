package migration

import (
	"car-rental-application/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Fungsi migrasi untuk create_cars_table
func createCarsTableMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240828174849_create_cars_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&models.Car{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.Car{})
		},
	}
}
