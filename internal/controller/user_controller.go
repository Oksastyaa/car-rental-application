package controller

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/service"
	"car-rental-application/pkg"
	"errors"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
	TopUpBalance(c echo.Context) error
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (uc *userController) RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "invalid input")
	}

	// validate user input
	if err := user.Validate(); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formattedErrors := pkg.FormatValidationError(&user, validationErrors)
			return pkg.RespondJSON(c, http.StatusBadRequest, nil, formattedErrors)
		}
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid Input"+err.Error())
	}

	if err := uc.userService.RegisterUser(&user); err != nil {
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to create user : "+err.Error())
	}

	return pkg.RespondJSON(c, http.StatusCreated, map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
	}, "User created successfully")
}

func (uc *userController) LoginUser(c echo.Context) error {
	user := new(models.LoginUser)
	if err := c.Bind(user); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid input")
	}
	// validate user input
	if err := validator.New().Struct(user); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formattedErrors := pkg.FormatValidationError(user, validationErrors)
			return pkg.RespondJSON(c, http.StatusBadRequest, nil, formattedErrors)
		}
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid Input: "+err.Error())
	}

	token, err := uc.userService.LoginUser(user.Email, user.Password)
	if err != nil {
		return pkg.RespondJSON(c, http.StatusUnauthorized, nil, err.Error())
	}
	return pkg.RespondJSON(c, http.StatusOK, map[string]any{
		"token": token,
	}, "Login success")

}

func (uc *userController) TopUpBalance(c echo.Context) error {
	userToken := c.Get("user")

	token, ok := userToken.(*jwt.Token)
	if !ok {
		return pkg.RespondJSON(c, http.StatusUnauthorized, nil, "Invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)

	userId := claims["id"].(float64)

	userIdUint := uint(userId)
	amount := struct {
		DepositAmount float64 `json:"deposit_amount" validate:"required,gt=0"`
	}{}

	if err := c.Bind(&amount); err != nil {
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid input")
	}

	if err := validator.New().Struct(amount); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			formattedErrors := pkg.FormatValidationError(amount, validationErrors)
			return pkg.RespondJSON(c, http.StatusBadRequest, nil, formattedErrors)
		}
		return pkg.RespondJSON(c, http.StatusBadRequest, nil, "Invalid Input: "+err.Error())
	}

	user, err := uc.userService.TopUpBalance(userIdUint, amount.DepositAmount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.RespondJSON(c, http.StatusNotFound, nil, "User not found")
		}
		return pkg.RespondJSON(c, http.StatusInternalServerError, nil, "Failed to top up balance")
	}

	return pkg.RespondJSON(c, http.StatusOK, map[string]any{
		"total_amount": user.DepositAmount,
	}, "Balance topped up successfully")
}
