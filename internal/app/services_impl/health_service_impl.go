package servicesimpl

import (
	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

type HealthServiceImpl struct {
	db database.HealthRepository
}

func NewHealthService(db database.HealthRepository) repos.HealthService {
	return &HealthServiceImpl{
		db: db,
	}
}

func (h *HealthServiceImpl) CheckServer() error {
	return h.db.PingDB()
}
