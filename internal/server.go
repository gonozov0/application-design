package internal

import (
	"log/slog"
	"net/http"

	_ "application-design/docs" // init swagger
	"application-design/internal/application/orders"
	hotelsInfra "application-design/internal/infrastructure/hotels"
	ordersInfra "application-design/internal/infrastructure/orders"
	usersInfra "application-design/internal/infrastructure/users"
	"application-design/pkg/echomiddleware"
	"application-design/pkg/environment"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

func newServer(config Config) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echomiddleware.SlogLoggerMiddleware(slog.Default()))
	e.Use(echomiddleware.PutRequestIDContext)
	e.Use(middleware.Recover())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(echomiddleware.PutSentryContext)

	if config.Server.Environment != environment.Production {
		e.GET("/swagger/*", swagger.WrapHandler)
		e.GET("/swagger", func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})
	}

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	userRepo := usersInfra.NewInMemoryRepo()
	hotelRepo := hotelsInfra.NewInMemoryRepo()
	orderRepo := ordersInfra.NewInMemoryRepo()
	orders.Setup(e, orderRepo, userRepo, hotelRepo)

	return e
}
