package service

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
	"fmt"
)

type TransactionService interface {
	UpdateTransactionStatus(InvoiceId, status, paymentMethod, paymentProvider string) error
	GetTransactionById(id uint) (*models.Transaction, error)
	GetAllTransaction() ([]models.Transaction, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) UpdateTransactionStatus(InvoiceId, status, paymentMethod, paymentProvider string) error {
	transaction, err := s.transactionRepo.GetTransactionByInvoiceID(InvoiceId)
	if err != nil {
		return fmt.Errorf("transaction Not Found: %w", err)
	}
	transaction.TransactionStatus = status
	transaction.PaymentMethod = paymentMethod
	transaction.PaymentProvider = paymentProvider

	if err := s.transactionRepo.UpdateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (s *transactionService) GetTransactionById(id uint) (*models.Transaction, error) {
	return s.transactionRepo.GetTransactionByID(id)
}

func (s *transactionService) GetAllTransaction() ([]models.Transaction, error) {
	return s.transactionRepo.GetAllTransactions()
}
