package repository

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
)

type RentalRepository interface {
	CreateRental(rental *models.Rental) (*models.Rental, error)
	GetRentalByID(id uint) (*models.Rental, error)
	GetAllRental() ([]models.Rental, error)
}

type rentalRepository struct {
	DB *gorm.DB
}

func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &rentalRepository{
		DB: db,
	}
}

func (r *rentalRepository) CreateRental(rental *models.Rental) (*models.Rental, error) {
	err := r.DB.Create(rental).Error
	if err != nil {
		return nil, err
	}
	return rental, nil
}
func (r *rentalRepository) GetRentalByID(id uint) (*models.Rental, error) {
	var rental models.Rental
	if err := r.DB.Preload("Transaction").First(&rental, id).Error; err != nil {
		return nil, err
	}
	return &rental, nil
}

func (r *rentalRepository) GetAllRental() ([]models.Rental, error) {
	var rentals []models.Rental
	if err := r.DB.Preload("Transaction").Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}
