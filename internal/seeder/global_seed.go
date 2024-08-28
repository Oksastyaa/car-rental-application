package seeder

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	SeedUsers(db)
	SeedCars(db)
	SeedRentals(db)
	SeedTransactions(db)
	logrus.Println("Seed all success")
}
