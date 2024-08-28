package migration

import (
	"car-rental-application/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Fungsi migrasi untuk create_users_table
func createUsersTableMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240828111020_create_users_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&models.User{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
