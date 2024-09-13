package postgres

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func (p *Postgres) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	var tenders []*models.Tender

	serviceTypes := *params.ServiceType

	// Преобразуем в строку и используем pq.Array дальше
	query := `
		SELECT id, name, description, service_type, status, organization_id, created_at
		FROM tenders
		WHERE service_type = ANY($1)
		LIMIT $2 OFFSET $3
	`

	rows, err := p.db.QueryContext(ctx, query, pq.Array(serviceTypes), params.Limit, params.Offset)
	if err != nil {
		p.logger.Error("Error in get tenders by service type", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
		if err != nil {
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		var currentVersion int32
		err = p.db.QueryRowContext(ctx, `
			SELECT COALESCE(MAX(version_number), 0)
			FROM tender_versions
			WHERE tender_id = $1`, tender.Id).Scan(&currentVersion)
		if err != nil {
			p.logger.Error("Error get current version number", zap.Error(err))
			return nil, err
		}
		tender.Version = currentVersion
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (p *Postgres) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	var tenders []*models.Tender
	var organizationId string

	id, err := p.GetUserIDByUsername(ctx, *params.Username)
	if id == -1 {
		p.logger.Error("User not found", zap.Error(err))
		return nil, ErrUserNotFound
	}
	if err != nil {
		p.logger.Error("Error in get user id by username", zap.Error(err))
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	query := `
        SELECT id, name, description, service_type, status, organization_id, created_at
        FROM tenders
        WHERE organization_id = $1
        LIMIT $2 OFFSET $3`

	rows, err := p.db.QueryContext(ctx, query, organizationId, params.Limit, params.Offset)
	if err != nil {
		p.logger.Error("Error get iorg tenders", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
		if err != nil {
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		var currentVersion int32
		err = p.db.QueryRowContext(ctx, `
			SELECT COALESCE(MAX(version_number), 0)
			FROM tender_versions
			WHERE tender_id = $1`, tender.Id).Scan(&currentVersion)
		if err != nil {
			p.logger.Error("Error get current version number", zap.Error(err))
			return nil, err
		}
		tender.Version = currentVersion
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (p *Postgres) CreateTender(ctx context.Context, params repos.CreateTenderParams) (*models.Tender, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error begin tx", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var organizationId string
	var tender models.Tender

	// Проверяем, что creator username является ответственным за организацию
	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.CreatorUsername).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	if *params.OrganizationID != organizationId {
		p.logger.Error("Not allowed")
		return nil, ErrUserNotAllowed
	}

	err = tx.QueryRowContext(ctx, `
    	INSERT INTO tenders (name, description, service_type, status, organization_id, created_at)
    	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
    	RETURNING id, name, description, service_type, status, organization_id, created_at`,
		params.Name, params.Description, params.ServiceType, params.Status, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		p.logger.Error("Error create tender", zap.Error(err))
		return nil, err
	}

	var version repos.TenderVersion
	err = tx.QueryRowContext(ctx, `
    	INSERT INTO tender_versions (tender_id, version_number, name, description, service_type, status, organization_id, created_at, is_current)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, TRUE)
    	RETURNING version_number`,
		tender.Id, 1, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId).Scan(&version)
	if err != nil {
		p.logger.Error("Error create new version of tender", zap.Error(err))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	tender.Version = version
	return &tender, nil
}

func (p *Postgres) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (*models.Tender, error) {
	var organizationId string
	var tender models.Tender

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error begin tx", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Валидируем пользователя: существует ли и является ли ответственным за организацию
	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	err = tx.QueryRowContext(ctx, `
        UPDATE tenders
        SET name = $1, description = $2, service_type = $3, updated_at = CURRENT_TIMESTAMP
        WHERE id = $4 AND organization_id = $5
        RETURNING id, name, description, service_type, status, organization_id, created_at`,
		params.Name, params.Description, params.ServiceType, tenderId, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		p.logger.Error("Error update tender", zap.Error(err))
		return nil, err
	}

	var currentVersion int32
	err = tx.QueryRowContext(ctx, `
        SELECT COALESCE(MAX(version_number), 0)
        FROM tender_versions
        WHERE tender_id = $1`, tenderId).Scan(&currentVersion)
	if err != nil {
		p.logger.Error("Error get current version number", zap.Error(err))
		return nil, err
	}

	tender.Version = currentVersion

	_, err = tx.ExecContext(ctx, `
        INSERT INTO tender_versions (tender_id, version_number, name, description, service_type, status, organization_id, created_at, updated_at, is_current)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, TRUE)`,
		tender.Id, tender.Version+1, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatedAt)
	if err != nil {
		p.logger.Error("Error create new version of tender", zap.Error(err))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	return &tender, nil
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (*models.Tender, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error begin tx", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var tender models.Tender
	var organizationId string

	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	var exists bool
	err = tx.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM tenders WHERE id = $1 AND organization_id = $2)`, tenderId, organizationId).Scan(&exists)
	if err != nil || !exists {
		p.logger.Error("Tender not found")
		return nil, ErrTenderNotFound
	}

	err = tx.QueryRowContext(ctx, `
        SELECT tender_id, name, description, service_type, status, organization_id, created_at
        FROM tender_versions
        WHERE tender_id = $1 AND version_number = $2`, tenderId, version).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		p.logger.Error("Version not found")
		return nil, ErrVersionNotFound
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE tender_versions
        SET is_current = FALSE
        WHERE tender_id = $1 AND is_current = TRUE`, tenderId)
	if err != nil {
		p.logger.Error("Error update is_current of version", zap.Error(err))
		return nil, err
	}

	var currentVersion int32
	err = tx.QueryRowContext(ctx, `
        SELECT COALESCE(MAX(version_number), 0)
        FROM tender_versions
        WHERE tender_id = $1`, tenderId).Scan(&currentVersion)
	if err != nil {
		p.logger.Error("Error get current version number", zap.Error(err))
		return nil, err
	}

	newVersion := currentVersion + 1
	_, err = tx.ExecContext(ctx, `
        INSERT INTO tender_versions (tender_id, version_number, name, description, service_type, status, organization_id, created_at, updated_at, is_current)
        VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, TRUE)`,
		tender.Id, newVersion, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId)
	if err != nil {
		p.logger.Error("Error insert new version in tender_versions", zap.Error(err))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	tender.Version = newVersion
	return &tender, nil
}

func (p *Postgres) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	var status repos.TenderStatus

	var organizationId string
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return "", ErrUserNotAllowed
	}

	err = p.db.QueryRowContext(ctx, `
        SELECT status
        FROM tenders
        WHERE id = $1 AND organization_id = $2`, tenderId, organizationId).Scan(&status)
	if err != nil {
		p.logger.Error("Tender not found")
		return "", ErrTenderNotFound
	}

	return status, nil
}

func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (*models.Tender, error) {
	var tender models.Tender

	var organizationId string
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	err = p.db.QueryRowContext(ctx, `
        UPDATE tenders
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2 AND organization_id = $3
        RETURNING id, name, description, service_type, status, organization_id, created_at`,
		params.Status, tenderId, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		p.logger.Error("Error update tender status", zap.Error(err))
		return nil, err
	}

	var currentVersion int32
	err = p.db.QueryRowContext(ctx, `
        SELECT COALESCE(MAX(version_number), 0)
        FROM tender_versions
        WHERE tender_id = $1`, tenderId).Scan(&currentVersion)
	if err != nil {
		p.logger.Error("Error get current version number", zap.Error(err))
		return nil, err
	}

	tender.Version = currentVersion
	return &tender, nil
}

func (p *Postgres) GetTenderByID(ctx context.Context, tenderId repos.TenderId) (*models.Tender, error) {
	var tender models.Tender

	err := p.db.QueryRowContext(ctx, `
        SELECT id, name, description, service_type, status, organization_id, created_at, version
        FROM tenders
        WHERE id = $1`, tenderId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		p.logger.Error("Tender not found")
		return nil, ErrTenderNotFound
	}

	return &tender, nil
}

func (p *Postgres) IsTenderExists(ctx context.Context, tenderId repos.TenderId) (bool, error) {
	var tenderExists bool
	tenderQuery := `SELECT EXISTS (SELECT 1 FROM tenders WHERE id = $1)`
	err := p.db.QueryRowContext(ctx, tenderQuery, tenderId).Scan(&tenderExists)
	if err != nil {
		return false, err
	}

	if !tenderExists {
		p.logger.Error("Tender not found")
		return false, ErrTenderNotFound
	}

	return true, nil
}
