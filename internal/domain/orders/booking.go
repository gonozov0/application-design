package orders

import (
	"errors"
	"time"

	"application-design/internal/domain"
	"application-design/internal/domain/internal"
)

var ErrInvalidBooking = errors.New("invalid booking")

type Booking struct {
	hotelID domain.HotelID
	roomID  domain.RoomID
	from    time.Time
	to      time.Time
}

func NewBooking(hotelID domain.HotelID, roomID domain.RoomID, from, to time.Time) (Booking, error) {
	if err := internal.ValidateBookingDates(from, to, ErrInvalidBooking); err != nil {
		return Booking{}, err
	}
	return Booking{
		hotelID: hotelID,
		roomID:  roomID,
		from:    from,
		to:      to,
	}, nil
}

func (b Booking) HotelID() domain.HotelID {
	return b.hotelID
}

func (b Booking) RoomID() domain.RoomID {
	return b.roomID
}

func (b Booking) From() time.Time {
	return b.from
}

func (b Booking) To() time.Time {
	return b.to
}
