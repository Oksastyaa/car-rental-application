package service

import (
	"car-rental-application/config"
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
	"car-rental-application/pkg"
	"context"
	"fmt"
	"github.com/xendit/xendit-go/v6/invoice"
	"time"
)

type RentalService interface {
	BookCar(rental *models.Rental) (*models.Rental, error)
	GetRentalByID(id uint) (*models.Rental, error)
	GetAllRental() ([]models.Rental, error)
}

type rentalService struct {
	rentalRepo      repository.RentalRepository
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepo
}

func NewRentalService(rentalRepo repository.RentalRepository, transactionRepo repository.TransactionRepository, userRepo repository.UserRepo) RentalService {
	return &rentalService{
		rentalRepo:      rentalRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *rentalService) BookCar(rental *models.Rental) (*models.Rental, error) {
	if rental.TotalCost <= 0 {
		return nil, nil
	}

	newRental, err := s.rentalRepo.CreateRental(rental)
	if err != nil {
		return nil, nil
	}
	if newRental.ID == 0 {
		return nil, nil
	}

	newRental.Transaction.RentalID = newRental.ID
	newRental.Transaction.UserID = newRental.UserID
	newRental.Transaction.InvoiceID = fmt.Sprintf("trx-%s", time.Now().Format("20060102150405"))
	newRental.Transaction.Amount = rental.TotalCost
	newRental.Transaction.TransactionStatus = "unpaid"
	newRental.Transaction.PaymentMethod = rental.Transaction.PaymentMethod
	newRental.Transaction.PaymentProvider = rental.Transaction.PaymentProvider
	newRental.Transaction.Description = fmt.Sprintf("Pembayaran rental mobil untuk %s", rental.RentalStartDate.Format("2006-01-02"))

	user, err := s.userRepo.FindByID(rental.UserID)
	if err != nil {
		return nil, nil
	}
	userEmail := user.Email
	invoiceID := fmt.Sprintf("trx-%s", time.Now().Format("20060102150405"))
	invoiceUrl, err := s.CreateXenditInvoice(invoiceID, userEmail, rental.TotalCost)
	if err != nil {
		return nil, nil
	}

	subject := "Konfirmasi Booking"
	emailBody := fmt.Sprintf("Booking Anda telah berhasil dibuat. Total yang harus dibayar adalah Rp%.2f. Silakan bayar melalui tautan berikut: %s", rental.TotalCost, invoiceUrl)
	err = pkg.SendEmail(userEmail, subject, emailBody)
	if err != nil {
		return nil, nil
	}

	return newRental, nil
}

func (s *rentalService) GetRentalByID(id uint) (*models.Rental, error) {
	return s.rentalRepo.GetRentalByID(id)
}

func (s *rentalService) GetAllRental() ([]models.Rental, error) {
	return s.rentalRepo.GetAllRental()
}

func (s *rentalService) CreateXenditInvoice(externalID, payerEmail string, amount float64) (string, error) {
	xenditClient := config.XenditClient
	description := "invoice for car rental"

	//using context background
	ctx := context.Background()

	resp, httpResponse, err := xenditClient.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(invoice.CreateInvoiceRequest{
			ExternalId:  externalID,
			Amount:      amount,
			PayerEmail:  &payerEmail,
			Description: &description,
		}).
		Execute()

	if err != nil {
		return "", fmt.Errorf("failed to create invoice: %w", err)
	}

	// check if the response status code is not 200
	if httpResponse.StatusCode != 200 {
		fmt.Printf("Failed to create invoice: %d\n", httpResponse.StatusCode)
		return "", fmt.Errorf("failed to create invoice: %s", httpResponse.Status)
	}

	return resp.InvoiceUrl, nil
}
