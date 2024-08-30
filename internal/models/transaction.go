package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	UserID            uint      `gorm:"not null" json:"user_id" validate:"required"`
	RentalID          uint      `gorm:"not null" json:"rental_id" validate:"required"`
	Amount            float64   `gorm:"not null" json:"amount" validate:"required,gt=0"`
	TransactionStatus string    `gorm:"type:varchar(50);not null" json:"transaction_status" validate:"required,oneof=paid unpaid refund"`
	TransactionDate   time.Time `gorm:"not null" json:"transaction_date" validate:"required"`
	InvoiceID         string    `gorm:"type:varchar(100);not null" json:"invoice_id" validate:"required"`
	PaymentMethod     string    `gorm:"type:varchar(50);not null" json:"payment_method" validate:"required,oneof=credit_card bank_transfer e_wallet"`
	PaymentProvider   string    `gorm:"type:varchar(50);not null" json:"payment_provider" validate:"required"`
	Description       string    `gorm:"type:text" json:"description"`
}

func (t *Transaction) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
