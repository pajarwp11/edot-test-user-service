package user

import (
	"context"
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
	ctx := context.Background()
	_, err := u.mysql.ExecContext(ctx, "INSERT INTO users (email,phone,password) VALUES (?,?,?)", user.Email, user.Phone, user.Password)
	return err
}

func (u *UserRepository) GetByEmailOrPhone(email string, phone string) (*user.User, error) {
	ctx := context.Background()
	user := user.User{}
	err := u.mysql.GetContext(ctx, &user, "SELECT email,phone,password FROM users WHERE email=? OR phone=?", email, phone)
	return &user, err
}
