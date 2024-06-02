package orders

import (
	"fmt"

	"application-design/internal/domain/hotels"
	"application-design/internal/domain/orders"
)

type Order struct {
	UserEmail string
	Bookings  []orders.Booking
}

func (s *OrderService) CreateOrder(dto Order) (*orders.Order, error) {
	user, err := s.userRepo.GetUser(dto.UserEmail)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	order, err := orders.CreateOrder(user.ID(), dto.Bookings)
	if err != nil {
		return nil, fmt.Errorf("could not create order: %w", err)
	}

	for _, booking := range dto.Bookings {
		hotel, err := s.hotelRepo.GetHotel(booking.HotelID())
		if err != nil {
			return nil, fmt.Errorf("could not get hotel: %w", err)
		}
		if err := hotel.BookRoom(hotels.BookingInfo{
			RoomID: booking.RoomID(),
			OderID: order.ID(),
			From:   booking.From(),
			To:     booking.To(),
		}); err != nil {
			return nil, fmt.Errorf("could not book room: %w", err)
		}
	}

	if err := s.orderRepo.SaveOrder(*order); err != nil {
		return nil, fmt.Errorf("could not save order: %w", err)
	}

	return order, nil
}
