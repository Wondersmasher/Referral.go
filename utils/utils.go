package utils

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Wondersmasher/Referral/env"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte((password + env.SALT)), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("could not hash password")
	}
	return string(hashed), nil
}

func ValidateHashedPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte((password+env.SALT))) == nil
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

type Claims struct {
	jwt.RegisteredClaims
	Email     string `json:"email"`
	UserName  string `json:"username"`
	CreatedAt string `json:"createdAt"`
}

func (c *Claims) NewClaims(duration time.Time) *Claims {

	return &Claims{
		Email:     c.Email,
		UserName:  c.UserName,
		CreatedAt: time.Now().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(duration),
			Issuer:    "referral-system-golang",
			Subject:   c.Email,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%s-%d", c.Email, time.Now().Unix()),
		},
	}
}
func CreateNewToken(email, userName string, duration time.Time, secretKey string) (string, error) {
	claims := Claims{
		Email:    email,
		UserName: userName,
	}
	time.Now()

	val, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.NewClaims(duration)).SignedString([]byte(secretKey))

	if err != nil {
		return "", errors.New("could not create token")

	}

	return val, nil

	// jwt.NewWithClaims(jwt.SigningMethodHS256, claims.NewClaims(duration)).SignedString([]byte("secret"))
}

func ValidateToken(tokenString string, secretKey string) (*Claims, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, false, fmt.Errorf("error parsing token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)

	if ok && token.Valid {
		return claims, true, nil
	} else {
		return nil, false, fmt.Errorf("invalid token")
	}
}
