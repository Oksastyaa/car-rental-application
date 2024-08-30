package service

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
)

type CarService interface {
	CreateCar(car *models.Car) (*models.Car, error)
	GetCarByID(id uint) (*models.Car, error)
	UpdateCar(car *models.Car) (*models.Car, error)
	DeleteCar(id uint) error
	GetAllCars() ([]models.Car, error)
}

type carService struct {
	carRepo repository.CarRepo
}

func NewCarService(carRepo repository.CarRepo) CarService {
	return &carService{
		carRepo: carRepo,
	}
}

func (cs *carService) CreateCar(car *models.Car) (*models.Car, error) {
	return cs.carRepo.CreateCar(car)
}

func (cs *carService) GetCarByID(id uint) (*models.Car, error) {
	return cs.carRepo.GetCarByID(id)
}

func (cs *carService) UpdateCar(car *models.Car) (*models.Car, error) {
	return cs.carRepo.UpdateCar(car)
}

func (cs *carService) DeleteCar(id uint) error {
	return cs.carRepo.DeleteCar(id)
}

func (cs *carService) GetAllCars() ([]models.Car, error) {
	return cs.carRepo.GetAllCars()
}
