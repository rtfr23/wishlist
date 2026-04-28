package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, conn_string string) (*pgxpool.Pool, error){
	return pgxpool.New(ctx, conn_string) 
}