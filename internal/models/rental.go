package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type Rental struct {
	gorm.Model
	UserID          uint        `gorm:"not null" json:"user_id" validate:"required"`
	CarID           uint        `gorm:"not null" json:"car_id" validate:"required"`
	RentalStartDate time.Time   `gorm:"not null" json:"rental_start_date" validate:"required"`
	RentalEndDate   *time.Time  `json:"rental_end_date"`
	TotalCost       float64     `gorm:"not null" json:"total_cost" validate:"required,gt=0"`
	Transaction     Transaction `gorm:"-" json:"transaction"`
}

func (r *Rental) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
