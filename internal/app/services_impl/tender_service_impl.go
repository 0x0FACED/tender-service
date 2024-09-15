package servicesimpl

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

// TODO: после вызова методов обращения к БД добавить проверку ошибки
// И вернуть ее
// В методах бд надо возвращать разные ошибки, чтобы потом определять, какой http код вернуть юзеру

type TenderServiceImpl struct {
	db database.TenderRepository
}

func NewTenderService(db database.TenderRepository) repos.TenderService {
	return &TenderServiceImpl{
		db: db,
	}
}

func (b *TenderServiceImpl) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	if err := validateGetTenders(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetTenders(ctx, params)
}

func (b *TenderServiceImpl) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	if err := validateGetUserTenders(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetUserTenders(ctx, params)
}

func (b *TenderServiceImpl) CreateTender(ctx context.Context, params repos.CreateTenderParams) (models.Tender, error) {
	if err := validateCreateTender(params); err != nil {
		return models.Tender{}, err.Error()
	}
	tender, err := b.db.CreateTender(ctx, params)
	if err != nil {
		return models.Tender{}, err
	}
	return *tender, nil
}

func (b *TenderServiceImpl) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (models.Tender, error) {
	if err := validateEditTender(params); err != nil {
		return models.Tender{}, err.Error()
	}
	tender, err := b.db.EditTender(ctx, tenderId, username, params)
	if err != nil {
		return models.Tender{}, err
	}
	return *tender, nil
}

func (b *TenderServiceImpl) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (models.Tender, error) {
	if err := validateRollbackTender(params); err != nil {
		return models.Tender{}, err.Error()
	}
	tender, err := b.db.RollbackTender(ctx, tenderId, version, params)
	if err != nil {
		return models.Tender{}, err
	}
	return *tender, nil
}

func (b *TenderServiceImpl) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	if err := validateGetTenderStatus(params); err != nil {
		return "", err.Error()
	}
	return b.db.GetTenderStatus(ctx, tenderId, params)
}

func (b *TenderServiceImpl) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (models.Tender, error) {
	if err := validateUpdateTenderStatus(params); err != nil {
		return models.Tender{}, err.Error()
	}
	tender, err := b.db.UpdateTenderStatus(ctx, tenderId, params)
	if err != nil {
		return models.Tender{}, err
	}
	return *tender, nil
}
