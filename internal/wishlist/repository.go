package wishlist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{database: db,}
}

func (r *Repository)InsertWishlist(ctx context.Context, wishlist Wishlist) (int, error) {
	sqlQuery := `
		INSERT INTO wishlists (user_id, event, description, date)
		VALUES($1, $2, $3, $4)
		RETURNING id;
	`
	var id int
	err := r.database.QueryRow(ctx, sqlQuery, wishlist.User_id,wishlist.Event, wishlist.Description, wishlist.Date).Scan(&id)
	return id, err	
}

func (r *Repository)SelectWishlist(ctx context.Context, wishlistId int, userId int) (Wishlist, error) {
	sqlQuery := `
		SELECT id, user_id, event, description, date
		FROM wishlists
		WHERE id = $1 AND user_id = $2;
	`
	wishlistRow := r.database.QueryRow(ctx, sqlQuery, wishlistId, userId)

	var wishlist Wishlist
	err := wishlistRow.Scan(&wishlist.Id, &wishlist.User_id, &wishlist.Event, &wishlist.Description, &wishlist.Date)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Wishlist{}, ErrWishlistNotFound
		} else {
			return Wishlist{}, err
		}
	}	

	return wishlist, nil
}

func (r *Repository)SelectAllWishlists(ctx context.Context, userId int) ([]Wishlist, error) {
	sqlQuery := `
		SELECT id, user_id, event, description, date
		FROM wishlists
		WHERE user_id = $1
		ORDER BY id ASC;
	`

	rows, err := r.database.Query(ctx, sqlQuery, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	wishlists := make([]Wishlist, 0)

	for rows.Next() {
		var wishlist Wishlist
		err := rows.Scan(&wishlist.Id, &wishlist.User_id, &wishlist.Event, &wishlist.Description, &wishlist.Date)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}

	return wishlists, nil
}

func (r *Repository)UpdateWishlist(ctx context.Context, wishlist Wishlist) (Wishlist, error) {
	sqlQuery := `
		UPDATE wishlists
		SET event = $1, description = $2, date = $3
		WHERE id = $4 AND user_id = $5
	`

	res, err := r.database.Exec(ctx, sqlQuery, wishlist.Event, wishlist.Description, wishlist.Date, wishlist.Id, wishlist.User_id)
	if err != nil {
		return Wishlist{}, err
	}	
	if res.RowsAffected() == 0 {
		return Wishlist{}, ErrWishlistNotFound
	}
	return wishlist, nil
}

func (r *Repository)DeleteWishlist(ctx context.Context, wishlistId int, userId int) error {
	sqlQuery := `
		DELETE FROM wishlists
		WHERE id = $1 AND user_id = $2;
	`
	res, err := r.database.Exec(ctx, sqlQuery, wishlistId, userId)

	if err != nil {
		return err
	}	
	
	if res.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}

	return nil
}

