package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string  `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password      string  `gorm:"not null" json:"password" validate:"required,min=6"`
	DepositAmount float64 `gorm:"default:0;not null" json:"deposit_amount" validate:"numeric"`
	Age           int     `gorm:"not null" json:"age" validate:"required"`
	Role          string  `gorm:"type:enum('admin','user');not null" json:"role" validate:"required,oneof=admin user"` // Only accept 'admin' or 'user'
	Token         string  `gorm:"type:varchar(255);" json:"token"`
}

// LoginUser is a struct to be used when logging in a user
type LoginUser struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
