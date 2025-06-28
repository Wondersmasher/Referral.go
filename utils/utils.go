package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"crypto/rand"
	"math/big"

	"reflect"
	"strconv"

	"github.com/Wondersmasher/Referral/env"
	"github.com/go-playground/validator/v10"
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

func ApiErrorResponse(err any) map[string]any {
	return map[string]any{
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
	Email      string `json:"email"`
	UserName   string `json:"username"`
	CreatedAt  string `json:"createdAt"`
	ReferralID string `json:"referralID"`
}

func (c *Claims) NewClaims(duration time.Time) *Claims {

	return &Claims{
		Email:      c.Email,
		UserName:   c.UserName,
		ReferralID: c.ReferralID,
		CreatedAt:  time.Now().String(),
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
func CreateNewToken(email, userName, referralID string, duration time.Time, secretKey string) (string, error) {
	claims := Claims{
		Email:    email,
		UserName: userName,
	}

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

func GenerateReferralID() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 6
	result := make([]byte, idLength)

	for i := 0; i < idLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return "REF-GO-" + strings.ToUpper(string(result)), nil
}

type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

func ValidateBodyRequest(payload interface{}) []*ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)

	if ok {
		reflected := reflect.ValueOf(payload)

		for _, validationErr := range validationErrors {
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}
			condition := validationErr.Tag()
			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			param := validationErr.Param()
			errMessage := keyToTitleCase + " field is " + condition

			switch condition {
			case "required":
				errMessage = keyToTitleCase + " is required"
			case "email":
				errMessage = keyToTitleCase + " must be a valid email address"
			case "min":
				if _, err := strconv.Atoi(param); err == nil {
					errMessage = fmt.Sprintf("%s must be at least %s characters", keyToTitleCase, param)
				}
			case "max":
				if _, err := strconv.Atoi(param); err == nil {
					errMessage = fmt.Sprintf("%s must be at ost %s characters", keyToTitleCase, param)
				}
			case "eqfield":
				errMessage = keyToTitleCase + " must be equal to " + strings.ToLower(param)
			}

			currentValidationError := &ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, currentValidationError)
		}
	}

	return errors
}
