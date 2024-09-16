/*
Package user provides HTTP handlers for user-related operations in the assessment application.

It includes functions for creating, retrieving, updating, and deleting user records in the database.
The package implements role-based access control, with certain operations restricted to admin users only.
It also handles request validation, password hashing, and JSON response formatting.
*/
package user

import (
	"assessment/dto"
	"assessment/internal/database"
	"assessment/internal/middleware"
	"assessment/models"
	"assessment/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func HandleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		middleware.AuthMiddleware(getUser)(w, r)
	case http.MethodPost:
		middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(createAdminUser))(w, r)
	case http.MethodPut:
		middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(updateUser))(w, r)
	case http.MethodDelete:
		middleware.AuthMiddleware(middleware.AdminOnlyMiddleware(deleteUser))(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// TODO : generate random password and email it to the email
func createAdminUser(w http.ResponseWriter, r *http.Request) {
	var adminUserRequest dto.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&adminUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	err = validate.Struct(adminUserRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Check if the email already exists
	var existingUser models.User
	result := database.DB.Where("email = ?", adminUserRequest.Email).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(adminUserRequest.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create the new admin user
	newAdminUser := models.User{
		Email:    adminUserRequest.Email,
		Password: hashedPassword,
		Role:     "admin", // Set the role to admin
	}

	result = database.DB.Create(&newAdminUser)
	if result.Error != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := dto.UserResponse{
		ID:    newAdminUser.ID,
		Email: newAdminUser.Email,
		Role:  newAdminUser.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

	var existingUser models.User
	result := database.DB.First(&existingUser, id)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the update request
	var updateRequest dto.UserUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validate.Struct(updateRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// check and update email
	if updateRequest.Email != "" && updateRequest.Email != existingUser.Email {
		var emailCheck models.User
		result := database.DB.Where("email = ?", updateRequest.Email).First(&emailCheck)
		if result.Error == nil {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}
		existingUser.Email = updateRequest.Email
	}

	// hash and update password
	if updateRequest.Password != "" {
		hashedPassword, err := utils.HashPassword(updateRequest.Password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		existingUser.Password = hashedPassword
	}

	// Update role
	if updateRequest.Role != "" {
		existingUser.Role = updateRequest.Role
	}

	result = database.DB.Save(&existingUser)
	if result.Error != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := dto.UserResponse{
		ID:    existingUser.ID,
		Email: existingUser.Email,
		Role:  existingUser.Role,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/user/")

	var existingUser models.User
	res := database.DB.First(&existingUser, id)
	if res.Error != nil {
		http.Error(w, "No user found", http.StatusNotFound)
		return
	}

	// Convert UserID from header to uint
	loggedInUserID, err := strconv.ParseUint(r.Header.Get("UserID"), 10, 32)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	if existingUser.ID == uint(loggedInUserID) {
		http.Error(w, "You cannot delete yourself !", http.StatusForbidden)
		return
	}

	result := database.DB.Unscoped().Delete(&models.User{}, id) // removing from db
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Fetch users with pagination
	var users []models.User
	var totalUsers int64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	result := database.DB.Offset(offset).Limit(pageSize).Find(&users)

	if result.Error != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	// Prepare pagination metadata
	totalPages := (int(totalUsers) + pageSize - 1) / pageSize
	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	response := struct {
		Users       []dto.UserResponse `json:"users"`
		TotalUsers  int64              `json:"totalUsers"`
		CurrentPage int                `json:"currentPage"`
		TotalPages  int                `json:"totalPages"`
		HasNextPage bool               `json:"hasNextPage"`
		HasPrevPage bool               `json:"hasPrevPage"`
	}{
		Users:       userResponses,
		TotalUsers:  totalUsers,
		CurrentPage: page,
		TotalPages:  totalPages,
		HasNextPage: hasNextPage,
		HasPrevPage: hasPrevPage,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
