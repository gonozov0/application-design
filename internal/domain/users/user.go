package users

import (
	"errors"
	"fmt"

	"application-design/internal/domain"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	id    domain.UserID
	email string
}

func NewUser(id domain.UserID, email string) (*User, error) {
	if !govalidator.IsEmail(email) {
		return nil, fmt.Errorf("%w: invalid email", ErrInvalidUser)
	}

	return &User{
		id:    id,
		email: email,
	}, nil
}

func CreateUser(email string) (*User, error) {
	return NewUser(domain.UserID(uuid.New()), email)
}

func (u *User) ID() domain.UserID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}
