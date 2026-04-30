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
	wishlistStr := mux.Vars(r)["id"]
	wishlistId, err := strconv.Atoi(wishlistStr)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

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

	item := Item{
		Wishlist_id: wishlistId,
	}

	if itemDto.Title != nil {
		item.Title = itemDto.Title
	}
	if itemDto.Description != nil {
		item.Description = itemDto.Description
	}
	if itemDto.URL != nil {
		item.URL = itemDto.URL
	}
	if itemDto.Priority != nil {
		item.Priority = itemDto.Priority
	}

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
		return
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

func (i *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	wishlistStr := mux.Vars(r)["id"]
	wishlistId, err := strconv.Atoi(wishlistStr)
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

	items, err := i.itemService.GetAllItems(r.Context(), wishlistId, userId)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	itemsDto := make([]dto.ItemDTO, 0, len(items))

	for _, a := range items {
		itemsDto = append(itemsDto, dto.ItemDTO{
			Id:          a.Id,
			Wishlist_id: a.Wishlist_id,
			Title:       a.Title,
			Description: a.Description,
			URL:         a.URL,
			Priority:    a.Priority,
		})
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(itemsDto, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}

func (i *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
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

	var patchedItemDto dto.ItemDTO
	if err := json.NewDecoder(r.Body).Decode(&patchedItemDto); err != nil {
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

	patchedItem := Item{
		Id:          itemId,
		Wishlist_id: wishlistId,
	}

	if patchedItemDto.Title != nil {
		patchedItem.Title = patchedItemDto.Title
	}

	if patchedItemDto.Description != nil {
		patchedItem.Description = patchedItemDto.Description
	}

	if patchedItemDto.URL != nil {
		patchedItem.URL = patchedItemDto.URL
	}
	if patchedItemDto.Priority != nil {
		patchedItem.Priority = patchedItemDto.Priority
	}
	item, err := i.itemService.UpdateItem(r.Context(), patchedItem, userId)

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

func (i *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
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

	err = i.itemService.DeleteItem(r.Context(), itemId, wishlistId, userId)

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

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if token == "" {
		errDTO := dto.NewErrorDTO(errors.New("token required"))
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

	err = h.itemService.ReserveItem(r.Context(), token, itemId)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrAlreadyReserved) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
