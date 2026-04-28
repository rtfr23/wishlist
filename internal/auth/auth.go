package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

type User struct {
	Email string
	Password string
}

func NewService(rep Repository) *Service {
	return &Service {
		repository: rep,
	}
}

func (s *Service)RegisterUser(ctx context.Context, user User) error {
	password_hash, err := generatePassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = password_hash
	s.repository.InsertUser(ctx, user)
	return nil
}

func generatePassword(password string) (string, error){
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(password_hash), nil
}


