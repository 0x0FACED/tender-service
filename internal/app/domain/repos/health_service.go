package repos

// HealthService предоставляет метод для проверки доступности сервера.
type HealthService interface {
	CheckServer() error
}
