package service

import (
	"car-rental-application/internal/models"
	"car-rental-application/internal/repository"
	"car-rental-application/pkg"
	"errors"
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
	user, err := u.userRepo.TopUpBalance(Id, amount)
	if err != nil {
		return nil, err
	}
	return user, nil
}
