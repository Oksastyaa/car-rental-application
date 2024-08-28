package service

import (
	"errors"
	"go-struktur-folder/internal/models"
	"go-struktur-folder/internal/repository"
	"go-struktur-folder/pkg"
)

type UserService interface {
	RegisterUser(user *models.User) error
	LoginUser(email string, password string) (string, error)
}

type userService struct {
	userRepo  repository.UserRepo
	jwtSecret string
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
