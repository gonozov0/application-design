package orders

import (
	"application-design/internal/application/orders/handlers"
	"application-design/internal/service/orders"

	"github.com/labstack/echo/v4"
)

func Setup(
	e *echo.Echo,
	orderRepo orders.OrderRepository,
	userRepo orders.UserRepository,
	hotelRepo orders.HotelRepository,
) {
	handler := handlers.NewHandler(orderRepo, userRepo, hotelRepo)

	orderGroup := e.Group("/orders")
	orderGroup.POST("", handler.Create)
}
