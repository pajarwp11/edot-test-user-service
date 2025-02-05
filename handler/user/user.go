package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"user-service/models/user"
)

type UserUsecase interface {
	Register(userRegister *user.RegisterRequest) error
	Login(userLogin *user.LoginRequest) (string, error)
}

type UserHandler struct {
	userUsecase UserUsecase
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewUserHandler(userUsecase UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, req *http.Request) {
	request := user.RegisterRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Email == "" || request.Phone == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "email, phone, and password are required"
		json.NewEncoder(w).Encode(response)
		return
	}

	err := u.userUsecase.Register(&request)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			w.WriteHeader(http.StatusBadRequest)
			response.Message = "email already registered"
		} else if strings.Contains(err.Error(), "phone") {
			w.WriteHeader(http.StatusBadRequest)
			response.Message = "phone number already registered"
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = "internal server error"
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "user registered"
	json.NewEncoder(w).Encode(response)
}

func (u *UserHandler) Login(w http.ResponseWriter, req *http.Request) {
	request := user.LoginRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Identification == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "email, phone, and password are required"
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := u.userUsecase.Login(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "login success"
	response.Data = map[string]interface{}{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}
