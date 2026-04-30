package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wishlist/internal/dto"
	"wishlist/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
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
		if errors.Is(err, ErrAccessDenied) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}

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

func (i *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	wishlistStr := mux.Vars(r)["id"]
	wishlistId, err := strconv.Atoi(wishlistStr)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	itemStr := mux.Vars(r)["itemid"]
	itemId, err := strconv.Atoi(itemStr)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
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

	item, err := i.itemService.GetItem(r.Context(), itemId, wishlistId, userId)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrItemNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	itemDto := dto.ItemDTO{
		Id:          item.Id,
		Wishlist_id: item.Wishlist_id,
		Title:       item.Title,
		Description: item.Description,
		URL:         item.URL,
		Priority:    item.Priority,
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(itemDto, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}
