package hotels

import (
	"sync"

	"application-design/internal/domain"
	"application-design/internal/domain/hotels"
)

type repoHotel struct {
	ID    domain.HotelID
	Rooms map[domain.RoomID]*repoRoom
}

type repoRoom struct {
	ID       domain.RoomID
	Bookings []hotels.Booking
}

type InMemoryRepo struct {
	hotels map[domain.HotelID]repoHotel
	mu     sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		hotels: make(map[domain.HotelID]repoHotel),
	}
}

func (r *InMemoryRepo) SaveHotel(h hotels.Hotel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hotels[h.ID()] = repoHotel{
		ID:    h.ID(),
		Rooms: make(map[domain.RoomID]*repoRoom),
	}
	return nil
}

func (r *InMemoryRepo) GetHotel(id domain.HotelID) (*hotels.Hotel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, ok := r.hotels[id]
	if !ok {
		return nil, hotels.ErrHotelNotFound
	}
	return hotels.NewHotel(h.ID, r), nil
}

func (r *InMemoryRepo) SaveRoom(room hotels.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hotel, ok := r.hotels[room.HotelID()]
	if !ok {
		return hotels.ErrHotelNotFound
	}

	existingRoom, ok := hotel.Rooms[room.ID()]
	if ok {
		existingRoom.ID = room.ID()
		return nil
	}
	hotel.Rooms[room.ID()] = &repoRoom{
		ID:       room.ID(),
		Bookings: make([]hotels.Booking, 0),
	}
	return nil
}

func (r *InMemoryRepo) GetRoom(hotelID domain.HotelID, id domain.RoomID) (*hotels.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hotel, ok := r.hotels[hotelID]
	if !ok {
		return nil, hotels.ErrHotelNotFound
	}
	room, ok := hotel.Rooms[id]
	if !ok {
		return nil, hotels.ErrRoomNotFound
	}
	return hotels.NewRoom(room.ID, hotel.ID, r), nil
}

func (r *InMemoryRepo) SaveBooking(hotelID domain.HotelID, booking hotels.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hotel, ok := r.hotels[hotelID]
	if !ok {
		return hotels.ErrHotelNotFound
	}
	room, ok := hotel.Rooms[booking.RoomID()]
	if !ok {
		return hotels.ErrRoomNotFound
	}
	for i, b := range room.Bookings {
		if b.OrderID() == booking.OrderID() {
			room.Bookings[i] = booking
			return nil
		}
	}
	room.Bookings = append(room.Bookings, booking)
	return nil
}

func (r *InMemoryRepo) GetBookings(
	hotelID domain.HotelID,
	orderID domain.OrderID,
	roomID domain.RoomID,
) ([]hotels.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hotel, ok := r.hotels[hotelID]
	if !ok {
		return nil, hotels.ErrHotelNotFound
	}
	room, ok := hotel.Rooms[roomID]
	if !ok {
		return nil, hotels.ErrRoomNotFound
	}
	bookings := make([]hotels.Booking, 0, len(room.Bookings))
	for _, b := range room.Bookings {
		if b.OrderID() == orderID {
			bookings = append(bookings, b)
		}
	}
	return bookings, nil
}

func (r *InMemoryRepo) IsBookingAvailable(hotelID domain.HotelID, booking hotels.Booking) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hotel, ok := r.hotels[hotelID]
	if !ok {
		return false
	}
	room, ok := hotel.Rooms[booking.RoomID()]
	if !ok {
		return false
	}
	for _, b := range room.Bookings {
		if booking.From().Before(b.To()) || booking.To().After(b.From()) {
			return false
		}
	}
	return true
}

func (r *InMemoryRepo) DeleteBooking(hotelID domain.HotelID, roomID domain.RoomID, id domain.BookingID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hotel, ok := r.hotels[hotelID]
	if !ok {
		return hotels.ErrHotelNotFound
	}
	room, ok := hotel.Rooms[roomID]
	if !ok {
		return hotels.ErrRoomNotFound
	}
	for i, b := range room.Bookings {
		if b.ID() == id {
			room.Bookings = append(room.Bookings[:i], room.Bookings[i+1:]...)
			return nil
		}
	}
	return hotels.ErrBookingNotFound
}
