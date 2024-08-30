package service

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xendit/xendit-go/v6/balance_and_transaction"
	"github.com/xendit/xendit-go/v6/common"
)

type HashPasswordFunc func(password string) (string, error)
type CheckPasswordHashFunc func(password, hash string) bool
type GenerateTokenFunc func(id uint, role, secret string) (string, error)

type UserServiceTest interface {
	RegisterUser(user *models.User) error
	LoginUser(email string, password string) (string, error)
	TopUpBalance(Id uint, amount float64) (*models.User, error)
}

type userServiceTest struct {
	userRepo          repository.UserRepo
	jwtSecret         string
	hashPassword      HashPasswordFunc
	checkPasswordHash CheckPasswordHashFunc
	generateToken     GenerateTokenFunc
	balanceApi        balance_and_transaction.BalanceApi
}

func NewUserServiceTest(userRepo repository.UserRepo, jwtSecret string, hashFunc HashPasswordFunc, checkHashFunc CheckPasswordHashFunc, generateTokenFunc GenerateTokenFunc, balanceApi balance_and_transaction.BalanceApi) UserServiceTest {
	return &userServiceTest{
		userRepo:          userRepo,
		jwtSecret:         jwtSecret,
		hashPassword:      hashFunc,
		checkPasswordHash: checkHashFunc,
		generateToken:     generateTokenFunc,
		balanceApi:        balanceApi,
	}
}

func (u *userServiceTest) RegisterUser(user *models.User) error {
	existingUser, _ := u.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashPassword, err := u.hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	return u.userRepo.CreateUser(user)
}

func (u *userServiceTest) LoginUser(email string, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !u.checkPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := u.generateToken(user.ID, user.Role, u.jwtSecret)
	if err != nil {
		return "", err
	}

	// Update token to user
	user.Token = token
	if err := u.userRepo.CreateUser(user); err != nil {
		return "", err
	}
	return token, nil
}

func (u *userServiceTest) TopUpBalance(Id uint, amount float64) (*models.User, error) {

	accountType := "CASH"
	currency := "IDR"

	resp, _, err := u.balanceApi.GetBalance(context.Background()).
		AccountType(accountType).
		Currency(currency).
		Execute()
	if err != nil {
		var xenditErr *common.XenditSdkError
		if errors.As(err, &xenditErr) {
			fmt.Printf("Error when calling `BalanceApi.GetBalance`: %v\n", xenditErr.Error())

			if xenditErr.FullError() != nil {
				fmt.Printf("Status Code: %v\n", xenditErr.Status())
			}
			fmt.Printf("Error Code: %v\n", xenditErr.Status())
			if xenditErr.FullError() != nil {
				fmt.Printf("Error Message: %v\n", xenditErr.Status())
			}

			return nil, fmt.Errorf("failed to get balance from Xendit: %v", xenditErr.FullError())
		}
	}
	currentBalance := resp.Balance
	newBalance := currentBalance + float32(amount)

	user, _ := u.userRepo.TopUpBalance(Id, float64(newBalance))

	return user, nil
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) TopUpBalance(id uint, newBalance float64) (*models.User, error) {
	args := m.Called(id, newBalance)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockBalanceApi struct {
	mock.Mock
}

// Implementasikan metode GetBalance yang mengembalikan objek request
func (m *MockBalanceApi) GetBalance(ctx context.Context) balance_and_transaction.ApiGetBalanceRequest {
	// Panggil metode mock untuk mencatat panggilan ini
	m.Called(ctx)
	// Kembalikan objek ApiGetBalanceRequest kosong untuk melanjutkan pengujian
	return balance_and_transaction.ApiGetBalanceRequest{}
}

func (m *MockBalanceApi) Execute() (*balance_and_transaction.Balance, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*balance_and_transaction.Balance), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBalanceApi) GetBalanceExecute(r balance_and_transaction.ApiGetBalanceRequest) (*balance_and_transaction.Balance, *http.Response, *common.XenditSdkError) {
	args := m.Called(r)
	if args.Get(0) != nil {
		return args.Get(0).(*balance_and_transaction.Balance), nil, nil
	}
	return nil, nil, args.Get(2).(*common.XenditSdkError)
}

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)

	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		func(password string) (string, error) {
			return "hashedPassword", nil
		},
		func(password, hash string) bool {
			return true
		},
		func(id uint, role, secret string) (string, error) {
			return "mockToken123", nil
		},
		new(MockBalanceApi),
	)

	user := &models.User{
		Email:    "Kartini_Wibisono@gmail.com",
		Password: "test12345",
	}

	mockRepo.On("FindByEmail", user.Email).Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	err := service.RegisterUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		nil,
		nil,
		nil,
		new(MockBalanceApi),
	)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &models.User{
		Email: "test@example.com",
	}

	// Setup mock behavior
	mockRepo.On("FindByEmail", user.Email).Return(existingUser, nil)

	// Call the service method
	err := service.RegisterUser(user)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_LoginUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)

	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		func(password string) (string, error) { return "hashedPassword", nil },
		func(password, hash string) bool { return true },
		func(id uint, role, secret string) (string, error) { return "token123", nil },
		new(MockBalanceApi),
	)

	user := &models.User{
		Email:    "Kartini_Wibisono@gmail.com",
		Password: "$2a$10$owV7isU9gMbzPS7euNcHoexNyk7qhEGxrzP15egnLXvppX5ZWDBXG",
		Role:     "user",
	}

	// Mock behavior
	mockRepo.On("FindByEmail", user.Email).Return(user, nil)
	mockRepo.On("CreateUser", mock.Anything).Return(nil)

	// Call the service method
	token, err := service.LoginUser(user.Email, "test12345")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)
	mockRepo.AssertExpectations(t)
}

func TestUserService_LoginUser_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		nil,
		nil,
		nil,
		new(MockBalanceApi),
	)

	// Setup mock behavior
	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	// Call the service method
	token, err := service.LoginUser("notfound@example.com", "password123")

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_TopUpBalance_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockBalanceApi := new(MockBalanceApi)

	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		func(password string) (string, error) { return "test12345", nil },
		func(password, hash string) bool { return true },
		func(id uint, role, secret string) (string, error) { return "token123", nil },
		mockBalanceApi,
	)

	user := &models.User{
		Email:         "test1@gmail.com",
		DepositAmount: 100,
	}

	balanceResponse := &balance_and_transaction.Balance{
		Balance: 0,
	}

	mockBalanceApi.On("GetBalance", mock.Anything).Return(balance_and_transaction.ApiGetBalanceRequest{})
	mockBalanceApi.On("Execute").Return(balanceResponse, nil, nil)
	mockRepo.On("TopUpBalance", uint(21), float64(150000)).Return(user, nil)

	result, err := service.TopUpBalance(21, 150000)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, float64(1000449936), result.DepositAmount)

	mockRepo.AssertExpectations(t)
	mockBalanceApi.AssertExpectations(t)
}
func TestUserService_TopUpBalance_Error(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockBalanceApi := new(MockBalanceApi)
	service := NewUserServiceTest(
		mockRepo,
		"jwtSecret",
		nil,
		nil,
		nil,
		mockBalanceApi,
	)

	// Setup mock behavior untuk metode GetBalance dan Execute
	mockBalanceApi.On("GetBalance", mock.Anything).Return(balance_and_transaction.ApiGetBalanceRequest{})
	mockBalanceApi.On("Execute").Return(nil, errors.New("failed to get balance from Xendit"))

	// Call the service method
	result, err := service.TopUpBalance(1, 500000)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get balance from Xendit")
	mockRepo.AssertExpectations(t)
}
