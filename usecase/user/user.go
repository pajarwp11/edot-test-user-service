package user

import (
	"errors"
	"time"
	"user-service/models/user"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Insert(user *user.User) error
	GetByEmailOrPhone(email string, phone string) (*user.User, error)
}

type UserUsecase struct {
	userRepo UserRepository
}

func NewUserUsecase(userRepo UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(userRegister *user.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userData := user.User{
		Email:    userRegister.Email,
		Phone:    userRegister.Phone,
		Password: string(hashedPassword),
	}
	err = u.userRepo.Insert(&userData)
	return err
}

func (u *UserUsecase) Login(userLogin *user.LoginRequest) (string, error) {
	user, err := u.userRepo.GetByEmailOrPhone(userLogin.Identification, userLogin.Identification)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	return generateJwt(user.Id)
}

func generateJwt(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString("myjwtsecret")
}
