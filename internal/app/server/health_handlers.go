package server

import "github.com/labstack/echo/v4"

// CheckServer converts echo context to params.
func (s *server) CheckServer(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = s.healthHandler.CheckServer(ctx)
	return err
}
