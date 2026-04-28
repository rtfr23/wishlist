package auth

import (
	"context"
	"time"
	"wishlist/internal/auth/token"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
	jwtMaker *token.JWTMaker
}

type User struct {
	Email string
	Password string
}

func NewService(rep *Repository, jwtM *token.JWTMaker) *Service {
	return &Service {
		repository: rep,
		jwtMaker: jwtM,
	}
}

func (s *Service)Login(ctx context.Context, user User) (string, error) {
	guessUser, err := s.repository.GetUser(ctx, user.Email)

	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(guessUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}

	token, err := s.jwtMaker.CreateToken(user.Email, 60*time.Minute)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service)Register(ctx context.Context, user User) error {
	password_hash, err := generatePassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = password_hash
	if err := s.repository.InsertUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func generatePassword(password string) (string, error){
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(password_hash), nil
}
