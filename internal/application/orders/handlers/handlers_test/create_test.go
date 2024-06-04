package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"application-design/internal/application/orders/handlers"
	"application-design/internal/domain"
	"application-design/internal/domain/hotels"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *OrderSuite) TestCreateOrderSuccess() {
	orderReq := handlers.CreateRequest{
		UserEmail: s.userEmail,
		HotelID:   s.hotel.ID(),
		RoomID:    s.roomID,
		From:      s.from,
		To:        s.to,
	}
	reqBody, _ := json.Marshal(orderReq)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusCreated, rec.Code)
	var resp handlers.CreateResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	orderID := resp.ID
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, orderID)

	bookings, err := s.hotelRepo.GetBookings(s.hotel.ID(), orderID, s.roomID)
	s.Require().NoError(err)
	s.Require().Len(bookings, 1)
	booking := bookings[0]
	s.Require().Equal(s.from, booking.From())
	s.Require().Equal(s.to, booking.To())
	_ = s.hotel.CancelBooking(s.roomID, booking.ID())
}

func (s *OrderSuite) TestCreateOrderInvalidDates() {
	tt := []struct {
		name   string
		from   string
		to     string
		errMsg string
	}{
		{
			name:   "from after to",
			from:   "4021-01-02T00:00:00Z",
			to:     "4021-01-01T00:00:00Z",
			errMsg: "invalid booking: \"from\" cannot be after \"to\"",
		},
		{
			name:   "from isn't a date",
			from:   "4021-01-01T01:00:00Z",
			to:     "4021-01-02T00:00:00Z",
			errMsg: "invalid booking: \"from\" must be a date without time",
		},
		{
			name:   "to isn't a date",
			from:   "4021-01-01T00:00:00Z",
			to:     "4021-01-02T01:00:00Z",
			errMsg: "invalid booking: \"to\" must be a date without time",
		},
		{
			name:   "from in the past",
			from:   "2020-01-01T00:00:00Z",
			to:     "4021-01-02T00:00:00Z",
			errMsg: "invalid booking: \"from\" cannot be in the past",
		},
		{
			name:   "to in the past",
			from:   "4021-01-01T00:00:00Z",
			to:     "2020-01-02T00:00:00Z",
			errMsg: "invalid booking: \"to\" cannot be in the past",
		},
	}

	for _, tc := range tt {
		s.Run(tc.name, func() {
			from, _ := time.Parse(time.RFC3339, tc.from)
			to, _ := time.Parse(time.RFC3339, tc.to)
			orderReq := handlers.CreateRequest{
				UserEmail: s.userEmail,
				HotelID:   s.hotel.ID(),
				RoomID:    s.roomID,
				From:      from,
				To:        to,
			}
			reqBody, _ := json.Marshal(orderReq)

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			s.e.ServeHTTP(rec, req)

			s.Require().Equal(http.StatusBadRequest, rec.Code)
			s.Require().Equal(tc.errMsg, rec.Body.String())
		})
	}
}

func (s *OrderSuite) TestCreateOrderUserNotFound() {
	orderReq := handlers.CreateRequest{
		UserEmail: "not_existed@user.com",
		HotelID:   s.hotel.ID(),
		RoomID:    s.roomID,
		From:      s.from,
		To:        s.to,
	}
	reqBody, _ := json.Marshal(orderReq)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusNotFound, rec.Code)
	s.Require().Equal("could not get user: user not found", rec.Body.String())
}

func (s *OrderSuite) TestCreateOrderHotelNotFound() {
	orderReq := handlers.CreateRequest{
		UserEmail: s.userEmail,
		HotelID:   "not_existed_hotel",
		RoomID:    s.roomID,
		From:      s.from,
		To:        s.to,
	}
	reqBody, _ := json.Marshal(orderReq)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusNotFound, rec.Code)
	s.Require().Equal("could not get hotel: hotel not found", rec.Body.String())
}

func (s *OrderSuite) TestCreateOrderRoomNotFound() {
	orderReq := handlers.CreateRequest{
		UserEmail: s.userEmail,
		HotelID:   s.hotel.ID(),
		RoomID:    "not_existed_room",
		From:      s.from,
		To:        s.to,
	}
	reqBody, _ := json.Marshal(orderReq)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusNotFound, rec.Code)
	s.Require().Equal("could not book room in service: could not get room: room not found", rec.Body.String())
}

func (s *OrderSuite) TestCreateOrderRoomAlreadyBooked() {
	bookingID, err := s.hotel.BookRoom(hotels.BookingInfo{
		RoomID: s.roomID,
		OderID: domain.OrderID(uuid.New()),
		From:   s.from,
		To:     s.to,
	})
	s.Require().NoError(err)
	defer func(roomID domain.RoomID, bookingID domain.BookingID) {
		_ = s.hotel.CancelBooking(roomID, bookingID)
	}(s.roomID, bookingID)

	orderReq := handlers.CreateRequest{
		UserEmail: s.userEmail,
		HotelID:   s.hotel.ID(),
		RoomID:    s.roomID,
		From:      s.from,
		To:        s.to,
	}
	reqBody, _ := json.Marshal(orderReq)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	s.e.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusConflict, rec.Code)
	s.Require().Equal("could not book room in service: could not book room: room already booked", rec.Body.String())
}
