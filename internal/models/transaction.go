package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	UserID          uint      `gorm:"not null" json:"user_id" validate:"required"`
	Amount          float64   `gorm:"not null" json:"amount" validate:"required,gt=0"`
	TransactionType string    `gorm:"type:varchar(50);not null" json:"transaction_type" validate:"required,oneof=top_up rent_payment refund"`
	TransactionDate time.Time `gorm:"not null" json:"transaction_date" validate:"required"`
}

func (t *Transaction) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
