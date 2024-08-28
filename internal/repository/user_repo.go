package repository

import (
	"go-struktur-folder/internal/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		DB: db,
	}
}

func (u *userRepo) CreateUser(user *models.User) error {
	return u.DB.Save(user).Error
}

func (u *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
