package item

import (
	"context"
)

type Service struct {
	repository *Repository
}

type Item struct {
	Id          int
	Wishlist_id int
	Title       string
	Description string
	URL         string
	Priority    int
	IsReserved  bool
}

func NewService(rep *Repository) *Service {
	return &Service{
		repository: rep,
	}
}

func (s *Service) AddItem(ctx context.Context, item Item, userId int) (int, error) {
	id, err := s.repository.InsertItem(ctx, item, userId)
	if err != nil {
		return 0, err

	}

	return id, nil
}

func (s *Service) GetItem(ctx context.Context, itemId int, wishlistId int, userId int) (Item, error) {
	item, err := s.repository.SelectItem(ctx, itemId, wishlistId, userId)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (s *Service) GetAllItems(ctx context.Context, userId int) ([]Item, error) {
	items, err := s.repository.SelectAllItems(ctx, userId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) UpdateItem(ctx context.Context, item Item, userId int) (Item, error) {
	updatedItem, err := s.repository.UpdateItem(ctx, item, userId)

	if err != nil {
		return Item{}, err
	}

	return updatedItem, nil
}

func (s *Service) DeleteItem(ctx context.Context, itemId int, userId int) error {
	if err := s.repository.DeleteItem(ctx, itemId, userId); err != nil {
		return err
	}

	return nil
}
