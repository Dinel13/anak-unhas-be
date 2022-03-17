package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dinel13/anak-unhas-be/helper"
)

// Checktoken check token for auth
func ChecToken(r *http.Request) (int, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return 0, errors.New("invalid token")
	}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	id, err := helper.ParseToken(tokenString)
	if err != nil {
		return 0, err

	}
	if id == 0 {
		return 0, errors.New("invalid token")

	}

	return id, nil
}

// CheckResetPasswordToken check token for reset password
func CheckResetPasswordToken(r *http.Request, jwtSecret string) (int, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return 0, errors.New("invalid token")
	}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	id, err := helper.ParseResetPasswordToken(tokenString)
	if err != nil {
		return 0, err

	}
	if id == 0 {
		return 0, errors.New("invalid token")

	}

	return id, nil
}
