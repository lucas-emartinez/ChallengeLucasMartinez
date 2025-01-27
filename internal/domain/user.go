package domain

import (
	"errors"
	"github.com/google/uuid"
)

// ErrUserNotFound es un error que se lanza cuando un usuario no se encuentra
var ErrUserNotFound = errors.New("user not found")

// User es una estructura que representa a un usuario
type User struct {
	ID        string
	Following []string
	Tweets    []Tweet
}

// NewUser crea una nueva instancia de User
func NewUser(username, email string) *User {
	return &User{
		ID: uuid.New().String(),
	}
}
