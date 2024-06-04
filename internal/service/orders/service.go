package orders

import (
	"application-design/internal/domain"
	"application-design/internal/domain/hotels"
	"application-design/internal/domain/orders"
	"application-design/internal/domain/users"
)

type OrderRepository interface {
	SaveOrder(order orders.Order) error
}

type UserRepository interface {
	GetUser(email string) (*users.User, error)
}

type HotelRepository interface {
	GetHotel(id domain.HotelID) (*hotels.Hotel, error)
	GetRoom(hotelID domain.HotelID, id domain.RoomID) (*hotels.Room, error)
	IsBookingAvailable(hotelID domain.HotelID, booking hotels.Booking) bool
	SaveBooking(hotelID domain.HotelID, booking hotels.Booking) error
}

type OrderService struct {
	orderRepo OrderRepository
	userRepo  UserRepository
	hotelRepo HotelRepository
}

func NewOrderService(or OrderRepository, ur UserRepository, hr HotelRepository) OrderService {
	return OrderService{
		orderRepo: or,
		userRepo:  ur,
		hotelRepo: hr,
	}
}
