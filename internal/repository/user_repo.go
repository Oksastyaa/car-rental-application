package repository

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	TopUpBalance(Id uint, newBalance float64) (*models.User, error)
	FindByID(Id uint) (*models.User, error)
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

func (u *userRepo) TopUpBalance(Id uint, newBalance float64) (*models.User, error) {
	var userToUpdate models.User
	err := u.DB.Where("id = ?", Id).First(&userToUpdate).Error
	if err != nil {
		return nil, err
	}
	userToUpdate.DepositAmount += newBalance

	err = u.DB.Save(&userToUpdate).Error
	if err != nil {
		return nil, err
	}
	return &userToUpdate, nil
}

func (u *userRepo) FindByID(Id uint) (*models.User, error) {
	var user models.User
	err := u.DB.Where("id = ?", Id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
