package user

type User struct {
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}
