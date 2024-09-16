package user

import (
	"assessment/dto"
	"assessment/internal/database"
	"assessment/internal/middleware"
	"assessment/models"
	"encoding/json"

	"net/http"
	"strings"
)

func HandleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		middleware.AuthMiddleware(getUser)(w, r)
	case http.MethodPut:
		middleware.AdminOnlyMiddleware(updateUser)
	case http.MethodDelete:
		middleware.AdminOnlyMiddleware(deleteUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

	var user models.User
	result := database.DB.First(&user, id)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	response := dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,   
		Role:     user.Role,    
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
