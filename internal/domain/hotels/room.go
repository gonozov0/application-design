package hotels

import (
	"errors"
	"fmt"
	"time"

	"application-design/internal/domain"
)

var (
	ErrRoomAlreadyBooked = errors.New("room already booked")
	ErrRoomNotFound      = errors.New("room not found")
)

type RoomRepository interface {
	IsBookingAvailable(hotelID domain.HotelID, booking Booking) bool
	SaveBooking(hotelID domain.HotelID, booking Booking) error
}

type Room struct {
	id      domain.RoomID
	hotelID domain.HotelID
	repo    RoomRepository
}

func NewRoom(id domain.RoomID, hotelID domain.HotelID, repo RoomRepository) *Room {
	return &Room{
		id:      id,
		hotelID: hotelID,
		repo:    repo,
	}
}

func (r *Room) ID() domain.RoomID {
	return r.id
}

func (r *Room) HotelID() domain.HotelID {
	return r.hotelID
}

func (r *Room) Book(orderID domain.OrderID, from, to time.Time) error {
	booking, err := NewBooking(orderID, r.id, from, to)
	if err != nil {
		return fmt.Errorf("could not create booking: %w", err)
	}
	if !r.repo.IsBookingAvailable(r.hotelID, *booking) {
		return ErrRoomAlreadyBooked
	}
	if err := r.repo.SaveBooking(r.hotelID, *booking); err != nil {
		return fmt.Errorf("could not save booking: %w", err)
	}
	return nil
}
