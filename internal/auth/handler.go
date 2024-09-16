package auth

import (
	"assessment/dto"
	"assessment/internal/database"
	"assessment/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

var validate = validator.New()

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginUser dto.UserRequest
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	result := database.DB.Where("email = ? AND password = ?", loginUser.Email, loginUser.Password).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// todo : implement user email check
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var UserRequest dto.UserRequest
	err := json.NewDecoder(r.Body).Decode(&UserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// request validation
	err = validate.Struct(UserRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	var user models.User
	res := database.DB.Where("email = ?", UserRequest.Email).First(&user)
	if res != nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	User := models.User{
		Email:    UserRequest.Email,
		Password: UserRequest.Password,
		Role:     UserRequest.Role,
	}

	result := database.DB.Create(&User)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(User) //case
}
