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
	panic("not implemented") // TODO: Implement
}

// Получение тендеров пользователя
func (b *TenderServiceImpl) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Создание нового тендера
func (b *TenderServiceImpl) CreateTender(ctx context.Context, params repos.CreateTenderParams) (models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Редактирование тендера
func (b *TenderServiceImpl) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Откат версии тендера
func (b *TenderServiceImpl) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Получение текущего статуса тендера
func (b *TenderServiceImpl) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	panic("not implemented") // TODO: Implement
}

// Изменение статуса тендера
func (b *TenderServiceImpl) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (models.Tender, error) {
	panic("not implemented") // TODO: Implement
}
