package user

import (
	"user-service/models/user"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	mysql *sqlx.DB
}

func NewUserRepository(mysql *sqlx.DB) *UserRepository {
	return &UserRepository{
		mysql: mysql,
	}
}

func (u *UserRepository) Insert(user *user.User) error {
	_, err := u.mysql.Exec("INSERT INTO users (email,phone,password) VALUES (?,?,?)", user.Email, user.Phone, user.Password)
	return err
}

func (u *UserRepository) GetByEmailOrPhone(email string, phone string) (*user.User, error) {
	user := user.User{}
	err := u.mysql.Get(&user, "SELECT id,email,phone,password FROM users WHERE email=? OR phone=?", email, phone)
	return &user, err
}
