package main

import (
	"context"
	"os"
	"wishlist/internal/auth"
	"wishlist/internal/db"
	"wishlist/internal/server"
)

func main(){
	ctx := context.Background()

	conn_string := os.Getenv("CONN_STRING")
	auth_pool, err := db.NewPostgresDB(ctx, conn_string)
	if err != nil {
		panic(err)
	}
	defer auth_pool.Close()

	auth_repository := auth.NewRepository(auth_pool)
	
	_ = auth_repository

	httpHandlers := server.NewHTTPHandlers()
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.Start(":8080")
	
}