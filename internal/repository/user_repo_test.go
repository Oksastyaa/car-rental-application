package repository

import (
	"car-rental-application/internal/models"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
)

var DB *gorm.DB

func setupTestDB() *gorm.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return DB
}

func Test_userRepo_CreateUser(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepo(db)

	user := &models.User{
		Email:         "mocktest@gmail.com",
		Password:      "password",
		DepositAmount: 0,
		Role:          "user",
		Age:           23,
	}

	// Test CreateUser
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	var createdUser models.User
	err = db.Where("email = ?", user.Email).First(&createdUser).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, createdUser.Email)

	// delete data after test
	db.Exec("DELETE FROM users WHERE email = ?", user.Email)
}

func Test_userRepo_FindByEmail_Success(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepo(db)

	user := &models.User{
		Email:    "Gilang_Kusuma@gmail.com",
		Password: "test12345",
	}
	db.Create(user)

	// call FindByEmail
	foundUser, err := repo.FindByEmail("Gilang_Kusuma@gmail.com")

	// asset not return error
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	// delete data after test
	db.Exec("DELETE FROM users WHERE email = ?", user.Email)
}

func TestFindByEmail_UserNotFound(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepo(db)

	foundUser, err := repo.FindByEmail("notfound@example.com")

	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestTopUpBalance_Success(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepo(db)

	user := &models.User{
		Email:         "test@example.com",
		DepositAmount: 100000,
	}
	db.Create(user)

	updatedUser, err := repo.TopUpBalance(user.ID, 50000)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, 150000.0, updatedUser.DepositAmount)

	db.Exec("DELETE FROM users WHERE email = ?", user.Email)
}

func TestTopUpBalance_UserNotFound(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepo(db)

	updatedUser, err := repo.TopUpBalance(999, 50000)

	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
