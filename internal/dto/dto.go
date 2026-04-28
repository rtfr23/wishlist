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