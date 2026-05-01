package dto

import (
	"encoding/json"
	"net/mail"
	"time"
)

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrDTO struct {
	Message string
	Time    time.Time
}

type WishlistDTO struct {
	Id          int        `json:"id"`
	Event       *string    `json:"event"`
	Description *string    `json:"desc"`
	Date        *time.Time `json:"date"`
	Token       string     `json:"token"`
}

type ItemDTO struct {
	Id          int     `json:"id"`
	Wishlist_id int     `json:"wishlist_id"`
	Title       *string `json:"title"`
	Description *string `json:"desc"`
	URL         *string `json:"url"`
	Priority    *int    `json:"priority"`
}

func (w *WishlistDTO) Validate() bool {
	if w.Event == nil || w.Date == nil {
		return false
	}
	if w.Description != nil && len(*w.Description) > 1024 {
		return false
	}

	if len(*w.Event) > 256 {
		return false
	}

	if w.Date.Before(time.Now()) {
		return false
	}
	return true
}

func (i *ItemDTO) Validate() bool {
	if i.Title == nil {
		return false
	}
	return true
}

func (u *UserDTO) Validate() bool {
	if u.Email == "" || u.Password == "" {
		return false
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return false
	}

	if len(u.Password) < 4 {
		return false
	}
	return true
}

func NewErrorDTO(err error) ErrDTO {
	return ErrDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e *ErrDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return err.Error()
	}

	return string(b)
}
