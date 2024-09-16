package dto

type UserRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
