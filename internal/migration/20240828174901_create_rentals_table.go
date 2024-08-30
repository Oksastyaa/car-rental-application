package migration

import (
	"car-rental-application/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Fungsi migrasi untuk create_rentals_table
func createRentalsTableMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240828174901_create_rentals_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&models.Rental{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.Rental{})
		},
	}
}
