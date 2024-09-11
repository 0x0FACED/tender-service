package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

var (
	ErrUserIsNotResponsible = errors.New("user is not responsible")
)

// Получение списка тендеров
func (p *Postgres) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	var tenders []*models.Tender

	// 1. Формируем запрос для получения тендеров по типу сервиса
	query := `
        SELECT id, name, description, service_type, status, organization_id, created_at, version
        FROM tenders
        WHERE service_type = ANY($1)
        LIMIT $2 OFFSET $3`

	rows, err := p.db.QueryContext(ctx, query, params.ServiceType, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Читаем результат и заполняем список тендеров
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (p *Postgres) GetUserTenders(ctx context.Context, params repos.GetUserTendersParams) ([]*models.Tender, error) {
	var tenders []*models.Tender
	var organizationId string

	// 1. Валидация пользователя: проверка, существует ли пользователь
	var userExists bool
	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM employee WHERE username = $1)`, params.Username).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("user not found")
	}

	// 2. Проверка, является ли пользователь ответственным за организацию
	err = p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		return nil, ErrUserIsNotResponsible
	}

	// 3. Получаем тендеры организации
	query := `
        SELECT id, name, description, service_type, status, organization_id, created_at, version
        FROM tenders
        WHERE organization_id = $1
        LIMIT $2 OFFSET $3`

	rows, err := p.db.QueryContext(ctx, query, organizationId, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 4. Читаем результат и заполняем список тендеров
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

// Создание нового тендера
func (p *Postgres) CreateTender(ctx context.Context, params repos.CreateTenderParams) (*models.Tender, error) {
	p.Connect()

	var organizationId string
	var tender models.Tender

	// 1. Проверяем, что CreatorUsername является ответственным за организацию
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.CreatorUsername).Scan(&organizationId)
	if err != nil {
		return nil, ErrUserIsNotResponsible
	}

	// 2. Создаем тендер
	err = p.db.QueryRowContext(ctx, `
    	INSERT INTO tenders (name, description, service_type, status, organization_id, created_at)
    	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
    	RETURNING id, name, description, service_type, status, organization_id, created_at`,
		params.Name, params.Description, params.ServiceType, params.Status, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		log.Println("Err: ", err)
		return nil, err
	}

	// 3. Возвращаем созданный тендер
	return &tender, nil
}

// Редактирование тендера
func (p *Postgres) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (*models.Tender, error) {
	var organizationId string
	var tender models.Tender

	// 1. Валидируем пользователя: существует ли и является ли ответственным за организацию
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, username).Scan(&organizationId)
	if err != nil {
		return nil, ErrUserIsNotResponsible
	}

	// 2. Обновляем тендер
	err = p.db.QueryRowContext(ctx, `
        UPDATE tenders
        SET name = $1, description = $2, service_type = $3, updated_at = CURRENT_TIMESTAMP
        WHERE id = $4 AND organization_id = $5
        RETURNING id, name, description, service_type, status, organization_id, created_at, version`,
		params.Name, params.Description, params.ServiceType, tenderId, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		return nil, err
	}

	// 3. Возвращаем обновленный тендер
	return &tender, nil
}

// Откат версии тендера
func (p *Postgres) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (*models.Tender, error) {
	var tender models.Tender

	// 1. Проверяем, существует ли пользователь и является ли он ответственным за организацию
	var organizationId string
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		return nil, ErrUserIsNotResponsible
	}

	// 2. Проверяем наличие тендера
	var exists bool
	err = p.db.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM tenders WHERE id = $1 AND organization_id = $2)`, tenderId, organizationId).Scan(&exists)
	if err != nil || !exists {
		return nil, ErrTenderNotFound
	}

	// 3. Получаем тендер нужной версии
	err = p.db.QueryRowContext(ctx, `
        SELECT id, name, description, service_type, status, organization_id, created_at, version
        FROM tender_versions
        WHERE id = $1 AND version = $2`, tenderId, version).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		return nil, fmt.Errorf("version not found for tender")
	}

	// 4. Обновляем основную запись тендера до откатной версии
	_, err = p.db.ExecContext(ctx, `
        UPDATE tenders
        SET name = $1, description = $2, service_type = $3, status = $4, version = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6`, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.Version, tenderId)
	if err != nil {
		return nil, err
	}

	return &tender, nil
}

// Получение текущего статуса тендера
func (p *Postgres) GetTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.GetTenderStatusParams) (repos.TenderStatus, error) {
	var status repos.TenderStatus

	// 1. Валидируем пользователя: является ли он ответственным за организацию
	var organizationId string
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		return "", ErrUserIsNotResponsible
	}

	// 2. Проверяем наличие тендера и получаем его статус
	err = p.db.QueryRowContext(ctx, `
        SELECT status
        FROM tenders
        WHERE id = $1 AND organization_id = $2`, tenderId, organizationId).Scan(&status)
	if err != nil {
		return "", ErrTenderNotFound
	}

	// 3. Возвращаем статус тендера
	return status, nil
}

// Изменение статуса тендера
func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderId repos.TenderId, params repos.UpdateTenderStatusParams) (*models.Tender, error) {
	var tender models.Tender

	// 1. Проверяем, является ли пользователь ответственным за организацию
	var organizationId string
	err := p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		return nil, ErrUserIsNotResponsible
	}

	// 2. Обновляем статус тендера
	err = p.db.QueryRowContext(ctx, `
        UPDATE tenders
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2 AND organization_id = $3
        RETURNING id, name, description, service_type, status, organization_id, created_at, version`,
		params.Status, tenderId, organizationId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		return nil, err
	}

	// 3. Возвращаем обновленный тендер
	return &tender, nil
}

func (p *Postgres) GetTenderByID(ctx context.Context, tenderId repos.TenderId) (*models.Tender, error) {
	var tender models.Tender

	// 1. Получаем тендер по ID
	err := p.db.QueryRowContext(ctx, `
        SELECT id, name, description, service_type, status, organization_id, created_at, version
        FROM tenders
        WHERE id = $1`, tenderId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
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
		return false, ErrTenderNotFound
	}

	return true, nil
}
