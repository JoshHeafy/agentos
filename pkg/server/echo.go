package server

import (
	"agentos/pkg/di"
	"agentos/pkg/timezone"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	ServiceName string
	Container   di.Container
	Echo        *echo.Echo
	PublicAPI   *echo.Group
}

func New(container di.Container, serviceName string) *Server {
	e := echo.New()

	// CORS restricted
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: container.Config.AllowedOrigins,
		AllowMethods: container.Config.AllowedMethods,
	}))

	publicGroup := e.Group("/v1/public")

	return &Server{
		ServiceName: serviceName,
		Container:   container,
		Echo:        e,
		PublicAPI:   publicGroup,
	}
}

func (s Server) Start() error {
	// Configures the timezone for the hole application
	loc, err := time.LoadLocation(timezone.AmericaLima)
	if err != nil {
		return err
	}

	time.Local = loc

	// Handle SIGINT (CTRL+C) gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	s.Echo.Server.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}

	srvErr := make(chan error, 1)

	// Start server
	go func() {
		srvErr <- s.Echo.Start(fmt.Sprintf(":%d", s.Container.Config.PortHTTP))
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	if err := s.Echo.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func (s Server) HealthCheckController(c echo.Context) error {
	if err := s.Container.DB.Ping(c.Request().Context()); err != nil {
		s.Container.Logger.Error(c.Request().Context(), "err when pinging databases", "error", err.Error(), "server_time", time.Now(), "service_name", s.ServiceName)

		return c.JSON(http.StatusInternalServerError, map[string]any{
			"status":       "error",
			"error":        err.Error(),
			"server_time":  time.Now(),
			"service_name": s.ServiceName,
			"query_params": c.QueryParams(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":       "ok",
		"server_time":  time.Now(),
		"service_name": s.ServiceName,
		"query_params": c.QueryParams(),
	})
}
