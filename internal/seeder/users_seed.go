package seeder

import (
	"car-rental-application/internal/models"
	"car-rental-application/pkg"
	"math/rand"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var users []models.User
	roles := []string{"admin", "user"}
	for i := 1; i <= 15; i++ {
		password, _ := pkg.HashPassword("test12345")
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
		logrus.Fatalf("Failed to seed users: %v", err)
	} else {
		logrus.Println("Seed users success")
	}
}
