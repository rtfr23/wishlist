package main

import (
	"context"
	"fmt"
	"os"
	"wishlist/internal/auth"
	"wishlist/internal/auth/token"
	"wishlist/internal/db"
	"wishlist/internal/item"
	"wishlist/internal/middleware"
	"wishlist/internal/server"
	"wishlist/internal/wishlist"
)

func main() {
	ctx := context.Background()

	conn_string := os.Getenv("CONN_STRING")
	dbpool, err := db.NewPostgresDB(ctx, conn_string)
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	secret_key := os.Getenv("SECRET_KEY")
	authJwtMaker := token.NewJWTMaker(secret_key)

	jwtMiddleware := middleware.NewJWTMiddleware(authJwtMaker)

	auth_repository := auth.NewRepository(dbpool)
	authService := auth.NewService(auth_repository, authJwtMaker)
	authHandler := auth.NewAuthHandler(authService)

	wishlist_repository := wishlist.NewRepository(dbpool)
	wishlistService := wishlist.NewService(wishlist_repository)
	wishlistHandler := wishlist.NewWishlistHandler(wishlistService)

	item_repository := item.NewRepository(dbpool)
	itemService := item.NewService(item_repository)
	itemHandler := item.NewItemHandler(itemService)
	httpHandlers := server.NewHTTPHandlers(*authHandler, *wishlistHandler, *itemHandler)
	httpServer := server.NewHTTPServer(httpHandlers)

	if err := httpServer.Start(":8080", *jwtMiddleware); err != nil {
		fmt.Println(err)
	}
}
