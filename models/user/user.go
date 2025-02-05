package user

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}

type UserLogin struct {
	Identification string `json:"identification"`
	Password       string `json:"password"`
}
