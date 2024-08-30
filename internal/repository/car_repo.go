package repository

import (
	"car-rental-application/internal/models"
	"gorm.io/gorm"
)

type CarRepo interface {
	CreateCar(car *models.Car) (*models.Car, error)
	GetCarByID(id uint) (*models.Car, error)
	UpdateCar(car *models.Car) (*models.Car, error)
	DeleteCar(id uint) error
	GetAllCars() ([]models.Car, error)
}

type carRepo struct {
	DB *gorm.DB
}

func NewCarRepo(db *gorm.DB) CarRepo {
	return &carRepo{
		DB: db,
	}
}

func (cr *carRepo) CreateCar(car *models.Car) (*models.Car, error) {
	if err := cr.DB.Create(car).Error; err != nil {
		return nil, err
	}
	return car, nil
}

func (cr *carRepo) GetCarByID(id uint) (*models.Car, error) {
	var car models.Car
	if err := cr.DB.First(&car, id).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (cr *carRepo) UpdateCar(car *models.Car) (*models.Car, error) {
	if err := cr.DB.Save(car).Error; err != nil {
		return nil, err
	}
	return car, nil
}

func (cr *carRepo) DeleteCar(id uint) error {
	if err := cr.DB.Delete(&models.Car{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (cr *carRepo) GetAllCars() ([]models.Car, error) {
	var cars []models.Car
	if err := cr.DB.Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}
