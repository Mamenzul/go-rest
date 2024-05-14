package types

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]User, error)
	CreateUser(User) error
	StoreResetToken(email string) (string, error)
	CheckResetToken(token string) (bool, error)
	UpdatePassword(email string, password string) error
	DeleteResetToken(token string) error
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type ResetPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordTokenPayload struct {
	Password string `json:"password" validate:"required,min=3,max=130"`
	Email    string `json:"email" validate:"required,email"`
}
