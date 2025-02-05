package main

import (
	"fmt"
	"log"
	"net/http"
	"user-service/db/mysql"
	userHandler "user-service/handler/user"
	userRepo "user-service/repository/user"
	userUsecase "user-service/usecase/user"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()

	router := mux.NewRouter()
	userRepository := userRepo.NewUserRepository(mysql.MySQL)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userHandler := userHandler.NewUserHandler(userUsecase)
	router.HandleFunc("user/register", userHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("user/login", userHandler.Login).Methods(http.MethodPost)

	fmt.Println("server is running")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
	}
}
