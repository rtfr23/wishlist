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

func (r *Repository)Insert_User(ctx context.Context, User UserModel){}

func (r *Repository)Select_User(ctx context.Context, User UserModel){}