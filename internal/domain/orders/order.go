package orders

import (
	"errors"
	"fmt"

	"application-design/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrder  = errors.New("invalid order")
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	id       domain.OrderID
	userID   domain.UserID
	bookings []Booking
}

func NewOrder(id domain.OrderID, userID domain.UserID, bookings []Booking) (*Order, error) {
	if len(bookings) == 0 {
		return nil, fmt.Errorf("%w: bookings cannot be empty", ErrInvalidOrder)
	}

	return &Order{
		id:       id,
		userID:   userID,
		bookings: bookings,
	}, nil
}

func CreateOrder(userID domain.UserID, bookings []Booking) (*Order, error) {
	id := domain.OrderID(uuid.New())
	return NewOrder(id, userID, bookings)
}

func (o *Order) ID() domain.OrderID {
	return o.id
}

func (o *Order) UserID() domain.UserID {
	return o.userID
}
