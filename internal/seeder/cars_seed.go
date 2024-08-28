package seeder

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
)

// SeedCars seeds the cars table with sample data
func SeedCars(db *gorm.DB) {
	var cars []models.Car
	carNames := []string{"Toyota Camry", "Honda Civic", "Ford Mustang", "Chevrolet Camaro", "BMW 3 Series"}
	categories := []string{"Sedan", "SUV", "Hatchback", "Convertible", "Coupe"}

	for i := 1; i <= 10; i++ {
		car := models.Car{
			Name:              carNames[rand.Intn(len(carNames))],
			StockAvailability: rand.Intn(10) + 1,
			RentalCost:        float64(rand.Intn(50000) + 100000),
			Category:          categories[rand.Intn(len(categories))],
		}
		cars = append(cars, car)
	}

	if err := db.Create(&cars).Error; err != nil {
		log.Fatalf("Failed to seed cars: %v", err)
	} else {
		log.Println("Seed cars success")
	}
}
