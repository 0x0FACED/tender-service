package servicesimpl

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/database"
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

// Получение списка предложений для тендера
func (t *TenderServiceImpl) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) error {
	panic("impl me")
}

// Просмотр отзывов на прошлые предложения
func (t *TenderServiceImpl) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) error {
	panic("impl me")
}

// Получение списка тендеров
func (t *TenderServiceImpl) GetTenders(ctx context.Context, params repos.GetTendersParams) error {
	panic("impl me")
}

// Получение тендеров пользователя
func (t *TenderServiceImpl) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) error {
	panic("impl me")
}

// Создание нового тендера
func (t *TenderServiceImpl) CreateTender(ctx context.Context, params repos.CreateTenderParams) error {
	panic("impl me")
}

// Редактирование тендера
func (t *TenderServiceImpl) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) error {
	panic("impl me")
}

// Откат версии тендера
func (t *TenderServiceImpl) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) error {
	panic("impl me")
}

// Получение текущего статуса тендера
func (t *TenderServiceImpl) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) error {
	panic("impl me")
}

// Изменение статуса тендера
func (t *TenderServiceImpl) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) error {
	panic("impl me")
}
