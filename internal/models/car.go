package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	Name              string  `gorm:"not null" json:"name" validate:"required"`
	StockAvailability int     `gorm:"not null" json:"stock_availability" validate:"gt=0"`
	RentalCost        float64 `gorm:"not null" json:"rental_cost" validate:"required,gt=0"`
	Brands            string  `gorm:"not null" json:"brands" validate:"required"`
}

func (c *Car) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
