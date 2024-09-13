package database

type HealthRepository interface {
	PingDB() error
}
