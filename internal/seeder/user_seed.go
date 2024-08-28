package seeder

import (
	"Rental-Mobil-App/internal/models"
	"Rental-Mobil-App/pkg"
	"math/rand"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var users []models.User

	for i := 1; i <= 15; i++ {
		password, _ := pkg.HashPassword("test12345")
		roles := []string{"admin", "user"}
		role := roles[rand.Intn(len(roles))]
		user := models.User{
			Email:         pkg.RandomEmail(),
			DepositAmount: float64(rand.Intn(1000000)),
			Age:           rand.Intn(50) + 18,
			Role:          role,
			Password:      password,
		}
		users = append(users, user)
	}
	if err := db.Create(&users).Error; err != nil {
		return
	}
	logrus.Println("Seed users success")
}
