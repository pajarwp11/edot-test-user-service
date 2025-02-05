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

func (u *UserUsecase) Register(user *user.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		return err
	}
	err = u.userRepo.Insert(user)
	return err
}

func (u *UserUsecase) Login(userLogin *user.UserLogin) (string, error) {
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
