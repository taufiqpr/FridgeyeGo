package helper

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTClaims(id int, name, email string, duration time.Duration) jwt.MapClaims {
	now := time.Now()
	exp := now.Add(duration)

	return jwt.MapClaims{
		"sub":   id,
		"name":  name,
		"email": email,
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
	}
}

func GenerateToken(secret string, id int, name, email string, duration time.Duration) (string, error) {
	claims := GenerateJWTClaims(id, name, email, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
