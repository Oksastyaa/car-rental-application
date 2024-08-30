package seeder

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// SeedTransactions seeds the transactions table with sample data
func SeedTransactions(db *gorm.DB) {
	var transactions []models.Transaction
	var rentals []models.Rental

	if err := db.Find(&rentals).Error; err != nil {
		log.Fatalf("Failed to fetch rentals: %v", err)
	}

	transactionTypes := []string{"paid", "unpaid", "refund"}
	paymentMethod := []string{"cash", "credit card", "debit card"}
	PaymentProvider := []string{"Xendit", "OVO", "GoPay"}

	for _, rental := range rentals {
		transaction := models.Transaction{
			UserID:            rental.UserID,
			Amount:            rental.TotalCost,
			TransactionStatus: transactionTypes[rand.Intn(len(transactionTypes))],
			TransactionDate:   time.Now().AddDate(0, 0, -rand.Intn(30)),
			InvoiceID:         "TRX-" + time.Now().Format("20060102") + randSeq(5),
			RentalID:          rental.ID,
			PaymentMethod:     paymentMethod[rand.Intn(len(paymentMethod))],
			PaymentProvider:   PaymentProvider[rand.Intn(len(PaymentProvider))],
			Description:       "Transaction for rental ID " + strconv.Itoa(int(rental.ID)),
		}

		transactions = append(transactions, transaction)
	}

	if err := db.Create(&transactions).Error; err != nil {
		log.Fatalf("Failed to seed transactions: %v", err)
	} else {
		log.Println("Seed transactions success")
	}
}

func randSeq(i int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, i)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
