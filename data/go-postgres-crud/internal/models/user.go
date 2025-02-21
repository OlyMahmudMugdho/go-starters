package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
