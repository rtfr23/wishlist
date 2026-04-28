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

func (s *Service)Login(ctx context.Context, user User) error {
	guessUser, err := s.repository.GetUser(ctx, user.Email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(guessUser.Password), []byte(user.Password)); err != nil {
		return err
	}
	return nil
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




