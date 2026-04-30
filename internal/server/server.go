package server

import (
	"errors"
	"net/http"
	"wishlist/internal/auth"
	"wishlist/internal/item"
	"wishlist/internal/middleware"
	"wishlist/internal/wishlist"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandler *HTTPHandlers
}

type HTTPHandlers struct {
	/*
		TODO:
			auth handler
			wishlist handler
			item handler
	*/

	AuthHandler     *auth.AuthHandler
	WishlistHandler *wishlist.WishlistHandler
	ItemHandler     *item.ItemHandler
}

func NewHTTPHandlers(AuthHandler auth.AuthHandler, WishlistHandler wishlist.WishlistHandler, ItemHandler item.ItemHandler) *HTTPHandlers {
	return &HTTPHandlers{
		AuthHandler:     &AuthHandler,
		WishlistHandler: &WishlistHandler,
		ItemHandler:     &ItemHandler,
	}
}

func NewHTTPServer(handler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandler: handler,
	}
}

func (s *HTTPServer) Start(port string, jwtMiddleware middleware.JWTMiddleware) error {
	router := mux.NewRouter()
	router.Path("/auth/signup").Methods("POST").HandlerFunc(s.httpHandler.AuthHandler.RegisterUser)
	router.Path("/auth/signin").Methods("POST").HandlerFunc(s.httpHandler.AuthHandler.LoginUser)

	router.Path("/public/wishlists/{token}").Methods("GET").HandlerFunc(s.httpHandler.WishlistHandler.GetWishlistWithToken)
	closed := router.PathPrefix("/wishlists").Subrouter()
	closed.Use(jwtMiddleware.Check)
	closed.Path("").Methods("POST").HandlerFunc(s.httpHandler.WishlistHandler.AddWishlist)
	closed.Path("").Methods("GET").HandlerFunc(s.httpHandler.WishlistHandler.GetWishlists)
	closed.Path("/{id}").Methods("GET").HandlerFunc(s.httpHandler.WishlistHandler.GetWishlist)
	closed.Path("/{id}").Methods("PATCH").HandlerFunc(s.httpHandler.WishlistHandler.UpdateWishlist)
	closed.Path("/{id}").Methods("DELETE").HandlerFunc(s.httpHandler.WishlistHandler.DeleteWishlist)

	closed.Path("/{id}/items").Methods("POST").HandlerFunc(s.httpHandler.ItemHandler.AddItem)
	closed.Path("/{id}/items").Methods("GET").HandlerFunc(s.httpHandler.ItemHandler.GetItems)
	closed.Path("/{id}/items/{itemid}").Methods("GET").HandlerFunc(s.httpHandler.ItemHandler.GetItem)
	closed.Path("/{id}/items/{itemid}").Methods("PACTH").HandlerFunc(s.httpHandler.ItemHandler.UpdateItem)
	closed.Path("/{id}/items/{itemid}").Methods("DELETE").HandlerFunc(s.httpHandler.ItemHandler.DeleteItem)

	if err := http.ListenAndServe(port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
