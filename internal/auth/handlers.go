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
	var new_user dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&new_user); err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}


	user := User {
		Email: new_user.Email,
		Password: new_user.Password,
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

