package repository

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) (*models.Transaction, error)
	GetTransactionByInvoiceID(invoiceID string) (*models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetAllTransactions() ([]models.Transaction, error)
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		DB: db,
	}
}

func (t *transactionRepository) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	err := t.DB.Create(transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *transactionRepository) GetTransactionByInvoiceID(invoiceID string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := t.DB.Where("invoice_id = ?", invoiceID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	err := t.DB.Save(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := t.DB.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := t.DB.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
