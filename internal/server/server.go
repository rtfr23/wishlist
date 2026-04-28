package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandler *HTTPHandlers
}

type HTTPHandlers struct {
	/*
	TODO:
		auth service
		wishlist service
		item service
	*/
}

func NewHTTPHandlers() *HTTPHandlers {
	return &HTTPHandlers{}
}

func NewHTTPServer(handler *HTTPHandlers) *HTTPServer{
	return &HTTPServer {
		httpHandler: handler,
	}
}

func (s *HTTPServer)Start(port string) error {
	router := mux.NewRouter()

	/*
	TODO:
	endpoints

	auth:
		POST /auth/signin
		POST /auth/signup
	
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