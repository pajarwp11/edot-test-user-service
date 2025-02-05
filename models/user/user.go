package user

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}

type LoginRequest struct {
	Identification string `json:"identification"`
	Password       string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
