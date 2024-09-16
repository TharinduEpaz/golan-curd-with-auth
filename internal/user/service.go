package user

import (
	"assessment/dto"
	"assessment/internal/database"
	"assessment/models"
)

func updateEmail(updateRequest dto.UserUpdateRequest) bool {
	var emailCheck models.User
		result := database.DB.Where("email = ?", updateRequest.Email).First(&emailCheck)
		if result.Error == nil {
			return false
		}
		return true
}