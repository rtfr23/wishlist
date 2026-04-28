package main

import (
	"context"
	"os"
	"wishlist/internal/auth"
	"wishlist/internal/auth/token"
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


	secret_key := os.Getenv("SECRET_KEY")
	authJwtMaker := token.NewJWTMaker(secret_key)

	//jwtMiddleware := middleware.NewJWTMiddleware(authJwtMaker)

	auth_repository := auth.NewRepository(auth_pool)
	authService := *auth.NewService(auth_repository, authJwtMaker)
	authHandler := auth.NewAuthHandler(ctx, &authService)

	httpHandlers := server.NewHTTPHandlers(*authHandler)
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.Start(":8080")
	
}