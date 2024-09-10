package postgres

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

// Получение списка тендеров
func (p *Postgres) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Получение тендеров пользователя
func (p *Postgres) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Создание нового тендера
func (p *Postgres) CreateTender(ctx context.Context, params repos.CreateTenderParams) (*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Редактирование тендера
func (p *Postgres) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Откат версии тендера
func (p *Postgres) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

// Получение текущего статуса тендера
func (p *Postgres) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	panic("not implemented") // TODO: Implement
}

// Изменение статуса тендера
func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) GetTenderByID(ctx context.Context, tenderId repos.TenderId) (*models.Tender, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) IsTenderExists(ctx context.Context, tenderId repos.TenderId) (bool, error) {
	var tenderExists bool
	tenderQuery := `SELECT EXISTS (SELECT 1 FROM tenders WHERE id = $1)`
	err := p.db.QueryRowContext(ctx, tenderQuery, tenderId).Scan(&tenderExists)
	if err != nil {
		return false, err
	}

	if !tenderExists {
		return false, ErrTenderNotFound
	}

	return true, nil
}
