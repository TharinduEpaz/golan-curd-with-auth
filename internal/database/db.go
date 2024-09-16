package database

import (
	"assessment/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// todo : load env
func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/golang-api?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&models.User{})
}