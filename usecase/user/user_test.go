package user

import (
	"errors"
	"testing"
	"user-service/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Insert(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmailOrPhone(email string, phone string) (*user.User, error) {
	args := m.Called(email, phone)
	return args.Get(0).(*user.User), args.Error(1)
}

// Test Register function
func TestRegister(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUsecase := NewUserUsecase(mockRepo)

	registerRequest := &user.RegisterRequest{
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	mockRepo.On("Insert", mock.AnythingOfType("*user.User")).Return(nil)

	err := userUsecase.Register(registerRequest)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test Login function (Success Case)
func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUsecase := NewUserUsecase(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &user.User{
		Id:       1,
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByEmailOrPhone", "test@example.com", "test@example.com").Return(testUser, nil)

	token, err := userUsecase.Login(&user.LoginRequest{Identification: "test@example.com", Password: "password123"})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

// Test Login function (Wrong Password)
func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUsecase := NewUserUsecase(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &user.User{
		Id:       1,
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByEmailOrPhone", "test@example.com", "test@example.com").Return(testUser, nil)

	token, err := userUsecase.Login(&user.LoginRequest{Identification: "test@example.com", Password: "wrongpassword"})
	assert.Error(t, err)
	assert.Equal(t, "wrong password", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// Test Login function (User Not Found)
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userUsecase := NewUserUsecase(mockRepo)

	mockRepo.On("GetByEmailOrPhone", "notfound@example.com", "notfound@example.com").Return((*user.User)(nil), errors.New("user not found"))

	token, err := userUsecase.Login(&user.LoginRequest{Identification: "notfound@example.com", Password: "password123"})
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
