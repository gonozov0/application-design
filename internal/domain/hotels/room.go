package hotels

import (
	"errors"
	"fmt"
	"time"

	"application-design/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrRoomAlreadyBooked = errors.New("room already booked")
	ErrRoomNotFound      = errors.New("room not found")
)

type RoomRepository interface {
	IsBookingAvailable(hotelID domain.HotelID, booking Booking) bool
	SaveBooking(hotelID domain.HotelID, booking Booking) error
	DeleteBooking(hotelID domain.HotelID, roomID domain.RoomID, id domain.BookingID) error
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

func (r *Room) Book(orderID domain.OrderID, from, to time.Time) (domain.BookingID, error) {
	empty := domain.BookingID(uuid.Nil)
	booking, err := CreateBooking(orderID, r.id, from, to)
	if err != nil {
		return empty, fmt.Errorf("could not create booking: %w", err)
	}
	if !r.repo.IsBookingAvailable(r.hotelID, *booking) {
		return empty, ErrRoomAlreadyBooked
	}
	if err := r.repo.SaveBooking(r.hotelID, *booking); err != nil {
		return empty, fmt.Errorf("could not save booking: %w", err)
	}
	return booking.ID(), nil
}

func (r *Room) CancelBooking(id domain.BookingID) error {
	return r.repo.DeleteBooking(r.hotelID, r.id, id)
}
