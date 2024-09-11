package servicesimpl

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

type TenderServiceImpl struct {
	db database.TenderRepository
}

func NewTenderService(db database.TenderRepository) repos.TenderService {
	return &TenderServiceImpl{
		db: db,
	}
}

// Получение списка тендеров
func (b *TenderServiceImpl) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	// TODO: Добавить валидацию!!!
	return b.db.GetTenders(ctx, params)
}

// Получение тендеров пользователя
func (b *TenderServiceImpl) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	// TODO: Добавить валидацию!!!
	return b.db.GetUserTenders(ctx, params)
}

// Создание нового тендера
func (b *TenderServiceImpl) CreateTender(ctx context.Context, params repos.CreateTenderParams) (models.Tender, error) {
	// TODO: Добавить валидацию!!!
	tender, err := b.db.CreateTender(ctx, params)
	return *tender, err
}

// Редактирование тендера
func (b *TenderServiceImpl) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (models.Tender, error) {
	// TODO: Добавить валидацию!!!
	tender, err := b.db.EditTender(ctx, tenderId, username, params)
	return *tender, err
}

// Откат версии тендера
func (b *TenderServiceImpl) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (models.Tender, error) {
	// TODO: Добавить валидацию!!!
	tender, err := b.db.RollbackTender(ctx, tenderId, version, params)
	return *tender, err
}

// Получение текущего статуса тендера
func (b *TenderServiceImpl) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	// TODO: Добавить валидацию!!!
	return b.db.GetTenderStatus(ctx, tenderId, params)
}

// Изменение статуса тендера
func (b *TenderServiceImpl) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (models.Tender, error) {
	// TODO: Добавить валидацию!!!
	tender, err := b.db.UpdateTenderStatus(ctx, tenderId, params)
	return *tender, err
}
