package database

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

type TenderRepository interface {
	// Получение списка тендеров
	GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error)
	// Получение тендеров пользователя
	GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error)
	// Создание нового тендера
	CreateTender(ctx context.Context, params repos.CreateTenderParams) (*models.Tender, error)
	// Редактирование тендера
	EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (*models.Tender, error)
	// Откат версии тендера
	RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (*models.Tender, error)
	// Получение текущего статуса тендера
	GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error)
	// Изменение статуса тендера
	UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (*models.Tender, error)
	GetTenderByID(ctx context.Context, tenderId repos.TenderId) (*models.Tender, error)

	IsTenderExists(ctx context.Context, tenderId repos.TenderId) (bool, error)
}
