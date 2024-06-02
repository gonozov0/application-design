package hotels

import (
	"errors"
	"time"

	"application-design/internal/domain"
)

var ErrInvalidBooking = errors.New("invalid reservation")

type Booking struct {
	orderID domain.OrderID
	roomID  domain.RoomID
	from    time.Time
	to      time.Time
}

func NewBooking(orderID domain.OrderID, roomID domain.RoomID, from, to time.Time) (*Booking, error) {
	if err := domain.ValidateBookingDates(from, to, ErrInvalidBooking); err != nil {
		return nil, err
	}
	return &Booking{
		orderID: orderID,
		roomID:  roomID,
		from:    from,
		to:      to,
	}, nil
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
