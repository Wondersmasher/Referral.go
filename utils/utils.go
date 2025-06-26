package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("could not hash password")
	}
	return string(hashed), nil
}

func ValidateHashedPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func CheckStrongPassword(password string) error {

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if !regexp.MustCompile("[0-9]").MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if !regexp.MustCompile("[a-z]").MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !regexp.MustCompile("[A-Z]").MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile("[^a-zA-Z0-9]").MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

func ApiErrorResponse(err string) map[string]string {
	return map[string]string{
		"error": err,
	}
}

func ApiSuccessResponse(data any, message string) map[string]any {
	return map[string]any{
		"data":    data,
		"message": message,
	}
}

func ApiSuccessWithPaginationResponse(data any, pagination any, message string) map[string]any {
	return map[string]any{
		"data":       data,
		"pagination": pagination,
		"message":    message,
	}
}
