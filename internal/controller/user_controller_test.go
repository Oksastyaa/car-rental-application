package controller

import (
	"car-rental-application/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) LoginUser(email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) TopUpBalance(Id uint, amount float64) (*models.User, error) {
	args := m.Called(Id, amount)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	controller := NewUserController(mockUserService)

	mockUserService.On("RegisterUser", mock.AnythingOfType("*models.User")).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"test@example.com","password":"password123","role":"user","age":23}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.RegisterUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "User created successfully")
	mockUserService.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	controller := NewUserController(mockUserService)

	// Mock data
	user := &models.LoginUser{
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserService.On("LoginUser", user.Email, user.Password).Return("token123", nil)

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute the function
	err := controller.LoginUser(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Login success")
	mockUserService.AssertExpectations(t)
}

func TestTopUpBalance_Success(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	controller := NewUserController(mockUserService)

	// Mock data
	user := &models.User{
		DepositAmount: 1500000,
	}
	mockUserService.On("TopUpBalance", uint(1), float64(500000)).Return(user, nil)

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/topup", strings.NewReader(`{"deposit_amount":500000}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate a valid JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = 1.0
	c.Set("user", token)

	// Execute the function
	err := controller.TopUpBalance(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Balance topped up successfully")
	mockUserService.AssertExpectations(t)
}
