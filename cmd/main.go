package main

import (
	"assessment/internal/auth"
	"assessment/internal/database"
	"assessment/internal/user"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()

	http.HandleFunc("/api/v1/login/", auth.HandleLogin)
	http.HandleFunc("/api/v1/register/", auth.HandleRegister)

	http.HandleFunc("/api/v1/user/", user.HandleUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
