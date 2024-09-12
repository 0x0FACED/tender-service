package postgres

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// Получение списка тендеров
func (p *Postgres) GetTenders(ctx context.Context, params repos.GetTendersParams) ([]*models.Tender, error) {
	var tenders []*models.Tender

	// Разыменование указателя и преобразование в нужный формат
	serviceTypes := *params.ServiceType // разыменовываем указатель

	// Преобразуем в строку, если нужно, и используем pq.Array
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
	// 2. Читаем результат и заполняем список тендеров
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

	// 1. Валидация пользователя: проверка, существует ли пользователь
	id, err := p.GetUserIDByUsername(ctx, *params.Username)
	if id == -1 {
		p.logger.Error("User not found", zap.Error(err))
		return nil, ErrUserNotFound
	}
	if err != nil {
		p.logger.Error("Error in get user id by username", zap.Error(err))
		return nil, err
	}

	// 2. Проверка, является ли пользователь ответственным за организацию
	err = p.db.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	// 3. Получаем тендеры организации
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

	// 4. Читаем результат и заполняем список тендеров
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

// Создание нового тендера
func (p *Postgres) CreateTender(ctx context.Context, params repos.CreateTenderParams) (*models.Tender, error) {
	// Начинаем транзакцию
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error begin tx", zap.Error(err))
		return nil, err
	}

	// Если что-то пойдет не так, откатываем транзакцию
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var organizationId string
	var tender models.Tender

	// 1. Проверяем, что CreatorUsername является ответственным за организацию
	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.CreatorUsername).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	// 2. Создаем тендер
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

	// 3. Записываем версию тендера в таблицу tender_versions
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

	// 4. Фиксируем транзакцию, если все прошло успешно
	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	// 5. Добавляем версию в объект тендера и возвращаем его
	tender.Version = version
	return &tender, nil
}

// Редактирование тендера
func (p *Postgres) EditTender(ctx context.Context, tenderId repos.TenderId, username repos.Username, params repos.EditTenderParams) (*models.Tender, error) {
	var organizationId string
	var tender models.Tender

	// Начинаем транзакцию
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error begin tx", zap.Error(err))
		return nil, err
	}

	// Если что-то пойдет не так, откатываем транзакцию
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Валидируем пользователя: существует ли и является ли ответственным за организацию
	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	// 2. Обновляем тендер
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

	// 3. Записываем новую версию тендера в таблицу tender_versions
	_, err = tx.ExecContext(ctx, `
        INSERT INTO tender_versions (tender_id, version_number, name, description, service_type, status, organization_id, created_at, updated_at, is_current)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, TRUE)`,
		tender.Id, tender.Version+1, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.CreatedAt)
	if err != nil {
		p.logger.Error("Error create new version of tender", zap.Error(err))
		return nil, err
	}

	// 4. Если все прошло успешно, фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	// 5. Возвращаем обновленный тендер
	return &tender, nil
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderId repos.TenderId, version int32, params repos.RollbackTenderParams) (*models.Tender, error) {
	// Начинаем транзакцию
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

	// 1. Проверяем, существует ли пользователь и является ли он ответственным за организацию
	err = tx.QueryRowContext(ctx, `
        SELECT r.organization_id
        FROM organization_responsible r
        JOIN employee e ON e.id = r.user_id
        WHERE e.username = $1`, params.Username).Scan(&organizationId)
	if err != nil {
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	// 2. Проверяем наличие тендера
	var exists bool
	err = tx.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM tenders WHERE id = $1 AND organization_id = $2)`, tenderId, organizationId).Scan(&exists)
	if err != nil || !exists {
		p.logger.Error("Tender not found")
		return nil, ErrTenderNotFound
	}

	// 3. Получаем тендер нужной версии из tender_versions
	err = tx.QueryRowContext(ctx, `
        SELECT tender_id, name, description, service_type, status, organization_id, created_at
        FROM tender_versions
        WHERE tender_id = $1 AND version_number = $2`, tenderId, version).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatedAt)
	if err != nil {
		p.logger.Error("Version not found")
		return nil, ErrVersionNotFound
	}

	// 4. Обновляем флаг is_current для текущей версии на false
	_, err = tx.ExecContext(ctx, `
        UPDATE tender_versions
        SET is_current = FALSE
        WHERE tender_id = $1 AND is_current = TRUE`, tenderId)
	if err != nil {
		p.logger.Error("Error update is_current of version", zap.Error(err))
		return nil, err
	}

	// 5. Получаем текущую версию, чтобы вычислить новую версию
	var currentVersion int32
	err = tx.QueryRowContext(ctx, `
        SELECT COALESCE(MAX(version_number), 0)
        FROM tender_versions
        WHERE tender_id = $1`, tenderId).Scan(&currentVersion)
	if err != nil {
		p.logger.Error("Error get current version number", zap.Error(err))
		return nil, err
	}

	// 6. Вставляем новую версию в tender_versions с номером на +1 от текущей версии
	newVersion := currentVersion + 1
	_, err = tx.ExecContext(ctx, `
        INSERT INTO tender_versions (tender_id, version_number, name, description, service_type, status, organization_id, created_at, updated_at, is_current)
        VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, TRUE)`,
		tender.Id, newVersion, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId)
	if err != nil {
		p.logger.Error("Error insert new version in tender_versions", zap.Error(err))
		return nil, err
	}

	// 7. Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	// 8. Возвращаем тендер с новой версией
	tender.Version = newVersion
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
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return "", ErrUserNotAllowed
	}

	// 2. Проверяем наличие тендера и получаем его статус
	err = p.db.QueryRowContext(ctx, `
        SELECT status
        FROM tenders
        WHERE id = $1 AND organization_id = $2`, tenderId, organizationId).Scan(&status)
	if err != nil {
		p.logger.Error("Tender not found")
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
		p.logger.Error("Error in check org responsible", zap.Error(err))
		return nil, ErrUserNotAllowed
	}

	// 2. Обновляем статус тендера
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
