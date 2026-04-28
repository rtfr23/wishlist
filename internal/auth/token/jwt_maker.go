package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey, }
}

func (m *JWTMaker) CreateToken(email string, duration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: email,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(m.secretKey))

	if err != nil {
		return "", err
	}

	return signed, nil
}
