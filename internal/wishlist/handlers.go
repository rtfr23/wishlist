package wishlist

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

type WishlistHandler struct {
	wishlistService *Service
}

func NewWishlistHandler(wishlistService *Service) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

func (wh *WishlistHandler) AddWishlist(w http.ResponseWriter, r *http.Request) {
	var wishlistDto dto.WishlistDTO
	if err := json.NewDecoder(r.Body).Decode(&wishlistDto); err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if !wishlistDto.Validate() {
		errDTO := dto.NewErrorDTO(errors.New("Invalid event name or date"))
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

	wishlist := Wishlist{
		User_id:     userId,
		Event:       wishlistDto.Event,
		Description: wishlistDto.Description,
		Date:        wishlistDto.Date,
	}

	id, token, err := wh.wishlistService.AddWishlist(r.Context(), wishlist)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	wishlist.Id = id
	wishlist.Token = token
	b, err := json.MarshalIndent(wishlist, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}

func (wh *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
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

	wishlist, err := wh.wishlistService.GetWishlist(r.Context(), wishlistId, userId)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrWishlistNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	wishlistDto := dto.WishlistDTO{
		Id:          wishlist.Id,
		Event:       wishlist.Event,
		Description: wishlist.Description,
		Date:        wishlist.Date,
		Token:       wishlist.Token,
	}

	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(wishlistDto, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}

}

func (wh *WishlistHandler) GetWishlists(w http.ResponseWriter, r *http.Request) {
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

	wishlists, err := wh.wishlistService.GetAllWishlists(r.Context(), userId)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	wishlistsDto := make([]dto.WishlistDTO, 0, len(wishlists))

	for _, a := range wishlists {
		wishlistsDto = append(wishlistsDto, dto.WishlistDTO{
			Id:          a.Id,
			Event:       a.Event,
			Description: a.Description,
			Date:        a.Date,
			Token:       a.Token,
		})
	}

	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(wishlistsDto, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}

}

func (wh *WishlistHandler) GetWishlistWithToken(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	if token == "" {
		errDTO := dto.NewErrorDTO(errors.New("Empty token"))
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	wishlist, err := wh.wishlistService.GetWishlistWithToken(r.Context(), token)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusNotFound)
		return
	}

	b, err := json.MarshalIndent(wishlist, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrWishlistNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}

	}

	if _, err := w.Write(b); err != nil {
		return
	}
}

func (wh *WishlistHandler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	wishlistStr := mux.Vars(r)["id"]
	wishlistId, err := strconv.Atoi(wishlistStr)
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	var patchedWishlistDto dto.WishlistDTO
	if err := json.NewDecoder(r.Body).Decode(&patchedWishlistDto); err != nil {
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

	patchedWishlist := Wishlist{
		Id:      wishlistId,
		User_id: userId,
	}

	if patchedWishlistDto.Event != nil {
		patchedWishlist.Event = patchedWishlistDto.Event
	}
	if patchedWishlistDto.Description != nil {
		patchedWishlist.Description = patchedWishlistDto.Description
	}
	if patchedWishlistDto.Date != nil {
		patchedWishlist.Date = patchedWishlistDto.Date
	}

	wishlist, err := wh.wishlistService.UpdateWishlist(r.Context(), patchedWishlist)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrWishlistNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	wishlistDto := dto.WishlistDTO{
		Id:          wishlist.Id,
		Event:       wishlist.Event,
		Description: wishlist.Description,
		Date:        wishlist.Date,
	}

	w.WriteHeader(http.StatusOK)

	b, err := json.MarshalIndent(wishlistDto, "", "\t")
	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}

}

func (wh *WishlistHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
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

	err = wh.wishlistService.DeleteWishlist(r.Context(), wishlistId, userId)

	if err != nil {
		errDTO := dto.NewErrorDTO(err)
		if errors.Is(err, ErrWishlistNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
