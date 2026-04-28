package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"wishlist/internal/auth/token"
	"wishlist/internal/dto"
)

type JWTMiddleware struct {
	maker *token.JWTMaker
}

var ClaimsKey = struct{}{}

func (m* JWTMiddleware)Check(h http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		header := r.Header.Get("Authorization")
		if header == "" {
			errDTO := dto.NewErrorDTO(errors.New("Missing authorization token"))
			http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
			return
		}

		headerInfo := strings.SplitN(header, " ", 2)
		if len(headerInfo) != 2 || headerInfo[0] != "Bearer" {
			errDTO := dto.NewErrorDTO(errors.New("Invalid format token"))
			http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
			return
		}

		tokenStr := headerInfo[1]
		claims, err := m.maker.VerifyToken(tokenStr)

		if err != nil {
			errDTO := dto.NewErrorDTO(errors.New("Invalid token"))
			http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		h(w, r.WithContext(ctx))
	}
} 
