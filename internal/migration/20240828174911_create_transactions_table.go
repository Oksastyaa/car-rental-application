package migration

import (
	"car-rental-application/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Fungsi migrasi untuk create_transactions_table
func createTransactionsTableMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240828174911_create_transactions_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&models.Transaction{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&models.Transaction{})
		},
	}
}
