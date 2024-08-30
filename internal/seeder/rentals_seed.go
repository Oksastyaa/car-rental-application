package seeder

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

// SeedRentals reseeds the rentals table with sample data
func SeedRentals(db *gorm.DB) {
	var rentals []models.Rental
	var users []models.User
	var cars []models.Car

	db.Find(&users)
	db.Find(&cars)

	for i := 1; i <= 15; i++ {
		user := users[rand.Intn(len(users))]
		car := cars[rand.Intn(len(cars))]

		startDate := time.Now().AddDate(0, 0, -rand.Intn(30))
		endDate := startDate.AddDate(0, 0, rand.Intn(5))

		rental := models.Rental{
			UserID:          user.ID,
			CarID:           car.ID,
			RentalStartDate: startDate,
			RentalEndDate:   &endDate,
			TotalCost:       float64(rand.Intn(500000) + 1000000),
		}

		if err := db.Create(&rental).Error; err != nil {
			log.Fatalf("Failed to create rental: %v", err)
		} else {
			log.Printf("Rental created successfully with ID: %d", rental.ID)
		}

		rentals = append(rentals, rental)
	}

	log.Println("Seed rentals success")
}
