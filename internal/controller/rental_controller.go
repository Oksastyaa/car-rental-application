package controller

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/service"
	"car-rental-application/pkg"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type RentalController interface {
	BookCar(c echo.Context) error
	GetRentalByID(c echo.Context) error
	GetAllRentals(c echo.Context) error
}

type rentalController struct {
	rentalService      service.RentalService
	transactionService service.TransactionService
}

func (r *rentalController) BookCar(c echo.Context) error {
	var rental models.Rental

	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return pkg.RespondJSON(c, http.StatusUnauthorized, nil, "Unauthorized: Invalid token")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return pkg.RespondJSON(c, http.StatusUnauthorized, nil, "Unauthorized: Invalid claims")
	}

	userId, ok := claims["id"].(float64)
	if !ok {
		return pkg.RespondJSON(c, http.StatusUnauthorized, nil, "Unauthorized: Invalid user ID")
	}

	if err := c.Bind(&rental); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Failed to bind request body: "+err.Error())
	}

	if rental.TotalCost <= 0 {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid rental data: Total cost must be greater than zero")
	}

	rental.UserID = uint(userId)

	bookedRental, err := r.rentalService.BookCar(&rental)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, bookedRental, "Failed to book car: "+err.Error())
	}

	if bookedRental == nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, bookedRental, "Failed to book car: Rental object is nil")
	}

	// Return response
	return pkg.RespondJSON(c, http.StatusCreated, map[string]any{
		"rental": bookedRental,
	}, "Car booked successfully")
}

func (r *rentalController) GetRentalByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "invalid rental id")
	}

	// not found
	if id == 0 {
		return pkg.RespondJSON(c, http.StatusNotFound, nil, "rental not found")
	}

	userToken := c.Get("user_id").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	rental, err := r.rentalService.GetRentalByID(uint(id))
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "failed to get rental by id: "+err.Error())
	}

	if rental.UserID != userId {
		return pkg.RespondJSON(c, http.StatusForbidden, nil, "access denied: not allowed to view this rental")
	}
	return pkg.RespondJSON(c, http.StatusOK, map[string]any{
		"rental":             rental,
		"transaction_status": rental.Transaction.TransactionStatus,
	}, "rental fetched successfully")
}

func (r *rentalController) GetAllRentals(c echo.Context) error {
	rentals, err := r.rentalService.GetAllRental()
	if err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to fetch rentals")
	}

	response := make([]map[string]any, len(rentals))
	for i, rental := range rentals {
		response[i] = map[string]any{
			"rental":            rental,
			"transactionStatus": rental.Transaction.TransactionStatus,
		}
	}
	return pkg.RespondJSON(c, http.StatusOK, rentals, "Rentals fetched successfully")
}

func NewRentalController(rentalService service.RentalService) RentalController {
	return &rentalController{
		rentalService: rentalService,
	}
}
