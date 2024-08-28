package controller

import (
	"errors"
	"go-struktur-folder/internal/models"
	"go-struktur-folder/internal/service"
	"go-struktur-folder/pkg"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
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
