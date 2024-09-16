/*
Package utils provides utility functions.

It includes functions for password hashing and comparison using the bcrypt algorithm.
These utilities are crucial for secure user authentication and password management
in the application.

This package relies on the golang.org/x/crypto/bcrypt package for cryptographic operations.
*/
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}
