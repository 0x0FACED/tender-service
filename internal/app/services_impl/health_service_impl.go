package servicesimpl

import (
	"github.com/labstack/echo/v4"
)

type HealthServiceImpl struct{}

func (h *HealthServiceImpl) CheckServer(ctx echo.Context) error { panic("impl me") }
