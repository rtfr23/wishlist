package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{database: db,}
}

func (r *Repository)InsertUser(ctx context.Context, User User){}

func (r *Repository)GetUser(ctx context.Context, User User){}