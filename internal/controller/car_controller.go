package controller

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/service"
	"car-rental-application/pkg"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CarController interface {
	CreateCar(c echo.Context) error
	GetCarByID(c echo.Context) error
	UpdateCar(c echo.Context) error
	DeleteCar(c echo.Context) error
	GetAllCars(c echo.Context) error
}

type carController struct {
	carService service.CarService
}

func NewCarController(carService service.CarService) CarController {
	return &carController{
		carService: carService,
	}
}

func (cc *carController) CreateCar(c echo.Context) error {
	car := new(models.Car)
	if err := c.Bind(car); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "invalid input")
	}

	// validate car input
	if err := car.Validate(); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formattedErrors := pkg.FormatValidationError(car, validationErrors)
			return pkg.RespondJSON(c, http.StatusBadRequest, nil, formattedErrors)
		}
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid input: "+err.Error())
	}

	createdCar, err := cc.carService.CreateCar(car)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to create car : "+err.Error())
	}
	return pkg.RespondJSON(c, http.StatusCreated, createdCar, "Car created successfully")
}

func (cc *carController) GetCarByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid ID")
	}

	// car not found
	if _, err := cc.carService.GetCarByID(uint(id)); err != nil {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "Car not found")
	}

	car, err := cc.carService.GetCarByID(uint(id))
	if err != nil {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "Car not found")
	}
	return pkg.RespondJSON(c, http.StatusOK, car, "Car found")
}

func (cc *carController) UpdateCar(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid ID")
	}

	car := new(models.Car)
	if err := c.Bind(car); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid input")
	}

	car.ID = uint(id)

	// car not found
	existingCar, err := cc.carService.GetCarByID(uint(id))
	if err != nil {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "Car not found")
	}

	// validate car input
	if err := car.Validate(); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formattedErrors := pkg.FormatValidationError(car, validationErrors)
			return pkg.RespondJSON(c, http.StatusBadRequest, nil, formattedErrors)
		}
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid input: "+err.Error())
	}

	//stock_availability
	if car.StockAvailability != 0 {
		car.StockAvailability += existingCar.StockAvailability
	} else {
		car.StockAvailability = existingCar.StockAvailability
	}

	updatedCar, err := cc.carService.UpdateCar(car)

	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to update car: "+err.Error())
	}
	return pkg.RespondJSON(c, http.StatusOK, updatedCar, "Car updated successfully")
}

func (cc *carController) DeleteCar(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid ID")
	}

	car, err := cc.carService.GetCarByID(uint(id))
	if err != nil {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "Car not found")
	}
	if err := cc.carService.DeleteCar(uint(id)); err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to delete car: "+err.Error())
	}

	return pkg.RespondJSON(c, http.StatusOK, map[string]any{
		"name":   car.Name,
		"brands": car.Brands,
	}, "Car deleted successfully")
}

func (cc *carController) GetAllCars(c echo.Context) error {
	cars, err := cc.carService.GetAllCars()
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to get cars: "+err.Error())
	}
	return pkg.RespondJSON(c, http.StatusOK, cars, "Cars found")
}
