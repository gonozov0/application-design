package handlers

import (
	"application-design/internal/service/orders"
)

type Handler struct {
	orderRepo orders.OrderRepository
	userRepo  orders.UserRepository
	hotelRepo orders.HotelRepository
}

func NewHandler(or orders.OrderRepository, ur orders.UserRepository, hr orders.HotelRepository) Handler {
	return Handler{
		orderRepo: or,
		userRepo:  ur,
		hotelRepo: hr,
	}
}
