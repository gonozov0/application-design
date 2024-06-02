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
	Bookings map[domain.OrderID][]hotels.Booking
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
		Bookings: make(map[domain.OrderID][]hotels.Booking),
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
	bookings := room.Bookings[booking.OrderID()]
	room.Bookings[booking.OrderID()] = append(bookings, booking)
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
	return room.Bookings[orderID], nil
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
	bookings := room.Bookings[booking.OrderID()]
	for _, b := range bookings {
		if booking.From().Before(b.To()) && booking.To().After(b.From()) {
			return false
		}
	}
	return true
}
