package main

import (
	"assessment/internal/auth"
	"assessment/internal/database"
	"assessment/internal/user"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var db *gorm.DB

func main() {
	database.InitDB()

	http.HandleFunc("/api/v1/login/", auth.HandleLogin)       //working
	http.HandleFunc("/api/v1/register/", auth.HandleRegister) //working

	http.HandleFunc("/api/v1/user/", user.HandleUser)

	// http.HandleFunc("/api/v1/user/", authMiddleware(adminOnlyMiddleware(handleAdminStats)))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/users/"):]
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/users/"):]
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleAdminStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userCount int64
	db.Model(&User{}).Count(&userCount)

	stats := map[string]interface{}{
		"total_users": userCount,
		"server_time": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
