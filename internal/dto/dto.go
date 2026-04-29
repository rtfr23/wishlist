package dto

import (
	"encoding/json"
	"time"
)

type UserDTO struct {
	Email string`json:"email"`
	Password string`json:"password"`
}

type ErrDTO struct {
	Message string 
	Time time.Time
}

type WishlistDTO struct {
	Id int`json:"id"`
	Event string`json:"event"`
	Description string`json:"desc"`
	Date *time.Time`date:"date"`
}

func(w *WishlistDTO)Validate() bool {
	if w.Event == "" || w.Date == nil{
		return false
	}
	return true
}

func (u *UserDTO)Validate() bool {
	if u.Email == "" || u.Password == "" {
		return false
	}

	return true
}

func NewErrorDTO(err error) ErrDTO {
	return ErrDTO {
		Message: err.Error(),
		Time: time.Now(),
	}
}

func (e *ErrDTO)ToString() string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return err.Error()
	}

	return string(b)
}
