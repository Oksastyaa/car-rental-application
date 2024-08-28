package seeder

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

// SeedTransactions seeds the transactions table with sample data
func SeedTransactions(db *gorm.DB) {
	var transactions []models.Transaction
	var users []models.User

	db.Find(&users)

	transactionTypes := []string{"top_up", "rent_payment", "refund"}

	for i := 1; i <= 20; i++ {
		user := users[rand.Intn(len(users))]
		transactionType := transactionTypes[rand.Intn(len(transactionTypes))]

		transaction := models.Transaction{
			UserID:          user.ID,
			Amount:          float64(rand.Intn(500000) + 100000),
			TransactionType: transactionType,
			TransactionDate: time.Now().AddDate(0, 0, -rand.Intn(30)),
		}
		transactions = append(transactions, transaction)
	}

	if err := db.Create(&transactions).Error; err != nil {
		log.Fatalf("Failed to seed transactions: %v", err)
	} else {
		log.Println("Seed transactions success")
	}
}
