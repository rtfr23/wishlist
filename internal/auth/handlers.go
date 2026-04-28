package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"wishlist/internal/dto"
)

type AuthHandler struct {
	authService *Service
}

func NewAuthHandler(ctx context.Context, authservice *Service) *AuthHandler {
	return &AuthHandler{
		authService: authservice,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request){
	var userDto dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if !userDto.Validate() {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	user := User {
		Email: userDto.Email,
		Password: userDto.Password,
	}

	if err := h.authService.Register(r.Context(), user); err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrUserAlreadyExists){
			http.Error(w, errDTO.ToString(), http.StatusConflict)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

	b, err := json.MarshalIndent(user.Email, "", "\t")
	if err != nil{
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request){
	var userDto dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if !userDto.Validate() {
		errDTO := dto.NewErrorDTO(errors.New("Invalid email or password"))
		http.Error(w, errDTO.ToString(),http.StatusBadRequest)
		return
	}

	user := User {
		Email: userDto.Email,
		Password: userDto.Password,
	}

	token, err := h.authService.Login(r.Context(), user)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrWrongPassword){
			http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}