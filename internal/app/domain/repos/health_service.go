package repos

import "github.com/labstack/echo/v4"

// HealthService предоставляет метод для проверки доступности сервера.
type HealthService interface {
	// Проверка доступности сервера
	CheckServer(ctx echo.Context) error
}
