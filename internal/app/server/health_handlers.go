package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CheckServer converts echo context to params.
func (s *server) CheckServer(ctx echo.Context) error {
	if err := s.healthHandler.CheckServer(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Database is not responding")
	}
	return ctx.JSON(http.StatusOK, "OK")
}
