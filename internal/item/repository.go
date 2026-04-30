package item

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
	return &Repository{database: db}
}

func (r *Repository) InsertItem(ctx context.Context, item Item) (int, error) {
	sqlQuery := `
		INSERT INTO items (wishlist_id, title, description, url, priority)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id int
	err := r.database.QueryRow(ctx, sqlQuery, item.Wishlist_id, item.Title, item.Description, item.URL, item.Priority).Scan(&id)
	return id, err
}

func (r *Repository) SelectItem(ctx context.Context, itemId int, wishlistId int) (Item, error) {
	sqlQuery := `
		SELECT id, wishlist_id, title, description, url, priority, is_reserved
		FROM items
		WHERE id = $1 AND wishlist_id = $2;
	`
	itemRow := r.database.QueryRow(ctx, sqlQuery, itemId, wishlistId)

	var item Item
	err := itemRow.Scan(&item.Id, &item.Wishlist_id, &item.Title, &item.Description, &item.URL, &item.Priority, &item.IsReserved)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Item{}, ErrItemNotFound
		} else {
			return Item{}, err
		}
	}

	return item, nil
}

func (r *Repository) SelectAllItems(ctx context.Context, wishlistId int) ([]Item, error) {
	sqlQuery := `
		SELECT id, wishlist_id, title, description, url, priority, is_reserved
		FROM items
		WHERE wishlist_id = $1
		ORDER BY id ASC;
	`

	rows, err := r.database.Query(ctx, sqlQuery, wishlistId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Wishlist_id, &item.Title, &item.Description, &item.URL, &item.Priority, &item.IsReserved)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *Repository) UpdateItem(ctx context.Context, item Item) (Item, error) {
	sqlQuery := `
		UPDATE items
		SET title = $1, description = $2, url = $3, priority = $4
		WHERE id = $5 AND wishlist_id = $6
	`

	res, err := r.database.Exec(ctx, sqlQuery, item.Title, item.Description, item.URL, item.Priority, item.Id, item.Wishlist_id)
	if err != nil {
		return Item{}, err
	}
	if res.RowsAffected() == 0 {
		return Item{}, ErrItemNotFound
	}
	return item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, itemId int, wishlist_id int) error {
	sqlQuery := `
		DELETE FROM items
		WHERE id = $1 AND wishlist_id = $2;
	`
	res, err := r.database.Exec(ctx, sqlQuery, itemId, wishlist_id)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrItemNotFound
	}

	return nil
}
