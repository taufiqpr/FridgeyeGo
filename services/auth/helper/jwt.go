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

// ParseToken validates a JWT and returns its claims
func ParseToken(secret, tokenString string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		if err == nil {
			return nil, jwt.ErrSignatureInvalid
		}
		return nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}
