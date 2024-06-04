package hotels

import (
	"errors"
	"time"

	"application-design/internal/domain"
	"application-design/internal/domain/internal"

	"github.com/google/uuid"
)

var (
	ErrInvalidBooking  = errors.New("invalid booking")
	ErrBookingNotFound = errors.New("booking not found")
)

type Booking struct {
	id      domain.BookingID
	orderID domain.OrderID
	roomID  domain.RoomID
	from    time.Time
	to      time.Time
}

func NewBooking(
	id domain.BookingID,
	orderID domain.OrderID,
	roomID domain.RoomID,
	from, to time.Time,
) (*Booking, error) {
	if err := internal.ValidateBookingDates(from, to, ErrInvalidBooking); err != nil {
		return nil, err
	}
	return &Booking{
		id:      id,
		orderID: orderID,
		roomID:  roomID,
		from:    from,
		to:      to,
	}, nil
}

func CreateBooking(orderID domain.OrderID, roomID domain.RoomID, from, to time.Time) (*Booking, error) {
	id := domain.BookingID(uuid.New())
	return NewBooking(id, orderID, roomID, from, to)
}

func (r *Booking) ID() domain.BookingID {
	return r.id
}

func (r *Booking) OrderID() domain.OrderID {
	return r.orderID
}

func (r *Booking) RoomID() domain.RoomID {
	return r.roomID
}

func (r *Booking) From() time.Time {
	return r.from
}

func (r *Booking) To() time.Time {
	return r.to
}
