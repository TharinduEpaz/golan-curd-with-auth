/*
Package middleware provides HTTP middleware functions for the assessment application.

It includes middleware for authentication using JWT tokens and role-based access control.
These middleware functions can be used to protect routes and ensure that only authorized
users can access certain parts of the application.

This package is designed to work with the net/http package and integrates with the
authentication system used in the application.
*/
package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value

		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		r.Header.Set("UserID", fmt.Sprintf("%v", claims["id"]))
		r.Header.Set("UserRole", fmt.Sprintf("%v", claims["role"]))
		next.ServeHTTP(w, r)
	}
}

func AdminOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("UserRole")
		if role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
