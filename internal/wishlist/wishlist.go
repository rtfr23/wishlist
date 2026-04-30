package wishlist

import (
	"context"
	"time"
)

type Service struct {
	repository *Repository
}

type Wishlist struct {
	Id          int
	User_id     int
	Event       *string
	Description *string
	Date        *time.Time
	Token       string
}

func NewService(rep *Repository) *Service {
	return &Service{
		repository: rep,
	}
}

func (s *Service) AddWishlist(ctx context.Context, wishlist Wishlist) (int, string, error) {
	id, token, err := s.repository.InsertWishlist(ctx, wishlist)
	if err != nil {
		return 0, "", err
	}

	return id, token, nil
}

func (s *Service) GetWishlist(ctx context.Context, wishlistId int, userId int) (Wishlist, error) {
	wishlist, err := s.repository.SelectWishlist(ctx, wishlistId, userId)
	if err != nil {
		return Wishlist{}, err
	}
	return wishlist, nil
}

func (s *Service) GetWishlistWithToken(ctx context.Context, token string) (Wishlist, error) {
	wishlist, err := s.repository.SelectWishlistWithToken(ctx, token)
	if err != nil {
		return Wishlist{}, err
	}
	return wishlist, nil
}

func (s *Service) GetAllWishlists(ctx context.Context, userId int) ([]Wishlist, error) {
	wishlists, err := s.repository.SelectAllWishlists(ctx, userId)
	if err != nil {
		return nil, err
	}
	return wishlists, nil
}

func (s *Service) UpdateWishlist(ctx context.Context, wishlist Wishlist) (Wishlist, error) {
	updatedWishlist, err := s.repository.UpdateWishlist(ctx, wishlist)

	if err != nil {
		return Wishlist{}, err
	}

	return updatedWishlist, nil
}

func (s *Service) DeleteWishlist(ctx context.Context, wishlistId int, userId int) error {
	err := s.repository.DeleteWishlist(ctx, wishlistId, userId)

	if err != nil {
		return err
	}

	return nil
}
