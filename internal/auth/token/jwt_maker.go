package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey: secretKey, }
}

func (m *JWTMaker) CreateToken(userId int, duration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: strconv.Itoa(userId),
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

func (m *JWTMaker) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if  !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}