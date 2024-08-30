package service

import (
	"car-rental-application/config"
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
	"car-rental-application/pkg"
	"context"
	"errors"
	"fmt"
	"github.com/xendit/xendit-go/v6/common"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(user *models.User) error
	LoginUser(email string, password string) (string, error)
	TopUpBalance(Id uint, amount float64) (*models.User, error)
}

type userService struct {
	userRepo  repository.UserRepo
	jwtSecret string
	DB        *gorm.DB
}

func NewUserService(userRepo repository.UserRepo, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userService) RegisterUser(user *models.User) error {
	existingUser, _ := u.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	//hash password
	hashPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	return u.userRepo.CreateUser(user)
}

func (u *userService) LoginUser(email string, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// check password
	if !pkg.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	//generate jwt token
	token, err := pkg.GenerateToken(user.ID, user.Role, u.jwtSecret)
	if err != nil {
		return "", err
	}

	//update token to user
	user.Token = token
	if err := u.userRepo.CreateUser(user); err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) TopUpBalance(Id uint, amount float64) (*models.User, error) {
	accountType := "CASH"

	currency := "IDR"

	resp, _, err := config.XenditClient.BalanceApi.GetBalance(context.Background()).
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
