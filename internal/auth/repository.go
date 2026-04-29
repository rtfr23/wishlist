package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{database: db,}
}

func (r *Repository)InsertUser(ctx context.Context, user User) error {
	sqlQuery := `
		INSERT INTO users (email, password_hash)
		VALUES($1, $2)
	`
	_, err := r.database.Exec(ctx, sqlQuery, user.Email, user.Password)
	if err != nil {
		var dbErr *pgconn.PgError

		if errors.As(err, &dbErr){
			if dbErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *Repository)GetUser(ctx context.Context, Email string) (User, error) {
	sqlQuery := `
		SELECT id, password_hash
		FROM users
		WHERE email = $1
	`
	UserRow := r.database.QueryRow(ctx, sqlQuery, Email)

	var user User
	err := UserRow.Scan(&user.Id, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		} else {
			return User{}, err
		}
	}	

	return user, nil
}