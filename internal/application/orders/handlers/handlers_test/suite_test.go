package handlers_test

import (
	"testing"
	"time"

	application "application-design/internal/application/orders"
	"application-design/internal/domain"
	"application-design/internal/domain/hotels"
	"application-design/internal/domain/users"
	hotelInfra "application-design/internal/infrastructure/hotels"
	ordersInfra "application-design/internal/infrastructure/orders"
	usersInfra "application-design/internal/infrastructure/users"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type OrderSuite struct {
	suite.Suite
	e         *echo.Echo
	userRepo  *usersInfra.InMemoryRepo
	hotelRepo *hotelInfra.InMemoryRepo
	orderRepo *ordersInfra.InMemoryRepo
	hotel     *hotels.Hotel
	roomID    domain.RoomID
	userEmail string
	from      time.Time
	to        time.Time
}

func (s *OrderSuite) SetupSuite() {
	s.e = echo.New()
	s.orderRepo = ordersInfra.NewInMemoryRepo()
	s.userRepo = usersInfra.NewInMemoryRepo()
	s.hotelRepo = hotelInfra.NewInMemoryRepo()
	application.Setup(s.e, s.orderRepo, s.userRepo, s.hotelRepo)

	s.userEmail = "test@test.com"
	user, _ := users.NewUser(domain.UserID(uuid.New()), s.userEmail)
	_ = s.userRepo.SaveUser(*user)

	s.hotel = hotels.NewHotel("reddison", s.hotelRepo)
	_ = s.hotelRepo.SaveHotel(*s.hotel)
	s.roomID = "lux"
	_ = s.hotel.AddRoom(s.roomID, s.hotelRepo)
	s.from, _ = time.Parse(time.RFC3339, "4021-01-01T00:00:00Z")
	s.to, _ = time.Parse(time.RFC3339, "4021-01-02T00:00:00Z")
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
