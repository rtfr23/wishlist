package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wishlist/internal/dto"
	"wishlist/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

type ItemHandler struct {
	itemService *Service
}

func NewItemHandler(itemService *Service) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

func (i *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var itemDto dto.ItemDTO
	if err := json.NewDecoder(r.Body).Decode(&itemDto); err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if !itemDto.Validate() {
		errDTO := dto.NewErrorDTO(errors.New("Invalid title"))
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.ClaimsKey).(*jwt.RegisteredClaims)
	if !ok {
		errDTO := dto.NewErrorDTO(errors.New("Unauthorized"))
		http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(claims.Subject)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
		return
	}

	item := Item{
		Wishlist_id: itemDto.Wishlist_id,
		Title:       itemDto.Title,
		Description: itemDto.Description,
		URL:         itemDto.URL,
		Priority:    itemDto.Priority,
	}

	id, err := i.itemService.AddItem(r.Context(), item, userId)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	item.Id = id
	b, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}
