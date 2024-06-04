package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"application-design/internal/domain"
	hotelsDomain "application-design/internal/domain/hotels"
	ordersDomain "application-design/internal/domain/orders"
	usersDomain "application-design/internal/domain/users"
	service "application-design/internal/service/orders"

	"github.com/labstack/echo/v4"
)

type CreateRequest struct {
	HotelID   domain.HotelID `json:"hotel_id"`
	RoomID    domain.RoomID  `json:"room_id"`
	UserEmail string         `json:"email"`
	From      time.Time      `json:"from"`
	To        time.Time      `json:"to"`
}

type CreateResponse struct {
	ID domain.OrderID `json:"id"`
}

// Create godoc
//
//	@Summary		Create and pay for an order
//	@Description	Creates an order based on the hotel and room IDs provided along with the booking period from the
//	request.
//	This process involves checking the availability of the room, verifying the user's existence,
//	and ensuring no overlapping bookings exist.
//
//	@Router			/orders [post]
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateRequest	true	"Order creation request"
//	@Success		201		{object}	CreateResponse	"Order successfully created with unique identifier"
//	@Failure		400		{object}	responses.Error	"Invalid booking dates or invalid data format"
//	@Failure		404		{object}	responses.Error	"Hotel, room, or user not found"
//	@Failure		409		{object}	responses.Error	"Room already booked for the selected dates"
//	@Failure		500		{object}	responses.Error	"Internal server error if order processing fails unexpectedly"
func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	booking, err := ordersDomain.NewBooking(req.HotelID, req.RoomID, req.From, req.To)
	if err != nil {
		if errors.Is(err, ordersDomain.ErrInvalidBooking) {
			// TODO: discuss about changing error contract
			return c.String(http.StatusBadRequest, err.Error())
		}
		slog.Error("failed to init booking", "err", err)
		return echo.ErrInternalServerError
	}

	ocs := service.NewOrderService(h.orderRepo, h.userRepo, h.hotelRepo)
	order, err := ocs.CreateOrder(service.Order{
		UserEmail: req.UserEmail,
		Bookings:  []ordersDomain.Booking{booking},
	})
	if err != nil {
		if errors.Is(err, usersDomain.ErrUserNotFound) || errors.Is(err, hotelsDomain.ErrHotelNotFound) ||
			errors.Is(err, hotelsDomain.ErrRoomNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, ordersDomain.ErrInvalidOrder) {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, hotelsDomain.ErrRoomAlreadyBooked) {
			return c.String(http.StatusConflict, err.Error())
		}
		slog.Error("failed to create order", "err", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, CreateResponse{ID: order.ID()})
}
