package hotels

import (
	"errors"
	"fmt"
	"time"

	"application-design/internal/domain"

	"github.com/google/uuid"
)

var ErrHotelNotFound = errors.New("hotel not found")

type HotelRepository interface {
	GetRoom(hotelID domain.HotelID, roomID domain.RoomID) (*Room, error)
	SaveRoom(room Room) error
}

type Hotel struct {
	id   domain.HotelID
	repo HotelRepository
}

func NewHotel(id domain.HotelID, repo HotelRepository) *Hotel {
	return &Hotel{
		id:   id,
		repo: repo,
	}
}

func (h *Hotel) ID() domain.HotelID {
	return h.id
}

func (h *Hotel) AddRoom(roomID domain.RoomID, repo RoomRepository) error {
	room := NewRoom(roomID, h.id, repo)
	if err := h.repo.SaveRoom(*room); err != nil {
		return fmt.Errorf("could not save room: %w", err)
	}
	return nil
}

type BookingInfo struct {
	RoomID domain.RoomID
	OderID domain.OrderID
	From   time.Time
	To     time.Time
}

func (h *Hotel) BookRoom(info BookingInfo) (domain.BookingID, error) {
	empty := domain.BookingID(uuid.Nil)
	room, err := h.repo.GetRoom(h.id, info.RoomID)
	if err != nil {
		return empty, fmt.Errorf("could not get room: %w", err)
	}
	bookingID, err := room.Book(info.OderID, info.From, info.To)
	if err != nil {
		return empty, fmt.Errorf("could not book room: %w", err)
	}
	return bookingID, nil
}

func (h *Hotel) CancelBooking(roomID domain.RoomID, bookingID domain.BookingID) error {
	room, err := h.repo.GetRoom(h.id, roomID)
	if err != nil {
		return fmt.Errorf("could not get room: %w", err)
	}
	if err := room.CancelBooking(bookingID); err != nil {
		return fmt.Errorf("could not cancel booking: %w", err)
	}
	return nil
}
