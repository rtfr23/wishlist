package server

import (
	"errors"
	"net/http"
	"wishlist/internal/auth"
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

	AuthHandler *auth.AuthHandler
	WishlistHandler *wishlist.WishlistHandler
}

func NewHTTPHandlers(AuthHandler auth.AuthHandler, WishlistHandler wishlist.WishlistHandler) *HTTPHandlers {
	return &HTTPHandlers{
		AuthHandler: &AuthHandler,
		WishlistHandler: &WishlistHandler,
	}
}

func NewHTTPServer(handler *HTTPHandlers) *HTTPServer{
	return &HTTPServer {
		httpHandler: handler,
	}
}

func (s *HTTPServer)Start(port string, jwtMiddleware middleware.JWTMiddleware) error {
	router := mux.NewRouter()
	router.Path("/auth/signup").Methods("POST").HandlerFunc(s.httpHandler.AuthHandler.RegisterUser)
	router.Path("/auth/signin").Methods("POST").HandlerFunc(s.httpHandler.AuthHandler.LoginUser)
	
	closed := router.PathPrefix("/wishlists").Subrouter()
	closed.Use(jwtMiddleware.Check)
	router.Path("").Methods("POST").HandlerFunc(s.httpHandler.WishlistHandler.AddWishlist)
	router.Path("").Methods("GET").HandlerFunc(s.httpHandler.WishlistHandler.GetWishlists)
	router.Path("/{id}").Methods("GET").HandlerFunc(s.httpHandler.WishlistHandler.GetWishlist)
	router.Path("/{id}").Methods("PATCH").HandlerFunc(s.httpHandler.WishlistHandler.UpdateWishlist)
	router.Path("/{id}").Methods("DELETE").HandlerFunc(s.httpHandler.WishlistHandler.DeleteWishlist)
	
	/*
	TODO:
	endpoints

	wishlist:
		closed:
			GET /wishlists
			GET /wishlists/{id}
			POST /wishlists
			PATCH /wishlists/{id}
			DELETE /wishlists/{id}
		opened:
			GET /wishlists/{token}

	item:
		closed:
			GET /wishlists/name/items/{id}
			GET /wishlists/name/items
			PATCH /wishlists/name/items/{id}
			POST /wishlists/name/items
			DELETE /wishlists/name/items/{id}
		opened:
			PATCH /wishlists/token/items/{id}
	*/

	if err := http.ListenAndServe(port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else {
			return err
		}
	}

	return nil
}