package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"

	_ "github.com/lib/pq"
)

var (
	ErrTenderNotFound       = errors.New("tender not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrBidNotFound          = errors.New("bid not found")
	ErrUserNotAllowed       = errors.New("user not allowed to view or update bid status")
)

// CreateBid создает новое предложение, если тендер и организация существуют
func (p *Postgres) CreateBid(ctx context.Context, params repos.CreateBidParams) (*models.Bid, error) {
	// Проверяем, существует ли тендер
	exists, err := p.IsTenderExists(ctx, *params.TenderID)
	if err != nil || !exists {
		return nil, err
	}

	// Если указан OrganizationID, проверяем существование организации
	if params.OrganizationID != nil {
		var orgExists bool
		orgQuery := `SELECT EXISTS (SELECT 1 FROM organization WHERE id = $1)`
		err := p.db.QueryRowContext(ctx, orgQuery, *params.OrganizationID).Scan(&orgExists)
		if err != nil {
			return nil, err
		}

		if !orgExists {
			return nil, ErrOrganizationNotFound
		}
	}

	// Начинаем транзакцию, чтобы сохранить данные в обе таблицы атомарно
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Откат транзакции при возникновении ошибки
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Создаем новое предложение (bid)
	bidQuery := `
		INSERT INTO bids (name, description, status, tender_id, author_type, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 
			(SELECT id FROM employee WHERE username = $6), 
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, description, status, tender_id, author_type, author_id, created_at, updated_at`

	row := tx.QueryRowContext(ctx, bidQuery,
		*params.Name,
		*params.Description,
		*params.Status,
		*params.TenderID,
		func() string {
			if params.OrganizationID != nil {
				return "Organization"
			}
			return "User"
		}(),
		*params.CreatorUsername,
	)

	bid := &models.Bid{}
	err = row.Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Создаем запись в таблице bid_versions для созданного предложения
	versionQuery := `
		INSERT INTO bid_versions (bid_id, version_number, author_id, status, created_at, is_current)
		VALUES ($1, $2, (SELECT id FROM employee WHERE username = $3), $4, CURRENT_TIMESTAMP, TRUE)`

	_, err = tx.ExecContext(ctx, versionQuery,
		bid.Id, // ID созданного предложения
		1,      // Первая версия
		*params.CreatorUsername,
		*params.Status,
	)
	if err != nil {
		return nil, err
	}

	// Если все прошло успешно, фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return bid, nil
}

// GetUserBids возвращает список предложений (bids) пользователя по его username с учетом пагинации
func (p *Postgres) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error) {
	// Проверяем, существует ли пользователь с указанным username
	var userID int
	userQuery := `SELECT id FROM employee WHERE username = $1`
	err := p.db.QueryRowContext(ctx, userQuery, *params.Username).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	// Формируем запрос для получения списка предложений пользователя
	bidQuery := `
		SELECT b.id, b.name, b.description, b.status, b.author_type, b.author_id, b.created_at, v.version_number
		FROM bids b
		JOIN bid_versions v ON b.id = v.bid_id
		WHERE b.author_id = $1 AND v.is_current = TRUE
		LIMIT $2 OFFSET $3`

	rows, err := p.db.QueryContext(ctx, bidQuery, userID, *params.Limit, *params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.Bid

	// Обрабатываем полученные строки
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.CreatedAt,
			&bid.Version,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	// Проверяем ошибки после итераций
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

// GetBidsForTender возвращает список предложений (bids) для указанного тендера с учетом пагинации
func (p *Postgres) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error) {
	// Проверяем, существует ли тендер с указанным tenderId
	exists, err := p.IsTenderExists(ctx, tenderId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTenderNotFound
	}

	// Проверяем, существует ли пользователь с указанным username
	var userID int
	userQuery := `SELECT id FROM employee WHERE username = $1`
	err = p.db.QueryRowContext(ctx, userQuery, params.Username).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	// Формируем запрос для получения списка предложений для данного тендера
	bidQuery := `
		SELECT b.id, b.name, b.description, b.status, b.author_type, b.author_id, b.created_at, v.version_number
		FROM bids b
		JOIN bid_versions v ON b.id = v.bid_id
		WHERE b.tender_id = $1 AND b.author_id = $2 AND v.is_current = TRUE
		LIMIT $3 OFFSET $4`

	rows, err := p.db.QueryContext(ctx, bidQuery, tenderId, userID, *params.Limit, *params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализируем срез для хранения результатов
	var bids []*models.Bid

	// Обрабатываем полученные строки
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.CreatedAt,
			&bid.Version,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	// Проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}

// GetBidStatus возвращает статус предложения по bidId
func (p *Postgres) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error) {
	// Проверяем, существует ли предложение с данным bidId
	var status repos.BidStatus
	query := `SELECT status FROM bids WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, bidId).Scan(&status)
	if err == sql.ErrNoRows {
		return "", ErrBidNotFound
	} else if err != nil {
		return "", err
	}

	// Проверяем, существует ли пользователь и является ли он ответственным за организацию тендера
	var orgId string
	userCheckQuery := `
		SELECT org.organization_id 
		FROM organization_responsible org 
		JOIN bids b ON org.organization_id = b.author_id
		JOIN employee e ON e.id = org.user_id
		WHERE e.username = $1 AND b.id = $2`

	err = p.db.QueryRowContext(ctx, userCheckQuery, params.Username, bidId).Scan(&orgId)
	if err == sql.ErrNoRows {
		return "", ErrUserNotAllowed
	} else if err != nil {
		return "", err
	}

	// 3. Возвращаем статус
	return status, nil
}

// UpdateBidStatus обновляет статус предложения по bidId
func (p *Postgres) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (*models.Bid, error) {
	// 1. Проверяем, существует ли предложение с данным bidId
	var currentStatus repos.BidStatus
	query := `SELECT status FROM bids WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, bidId).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return nil, ErrBidNotFound
	} else if err != nil {
		return nil, err
	}

	// 2. Проверяем, существует ли пользователь и является ли он ответственным за организацию тендера
	var orgId string
	userCheckQuery := `
		SELECT org.organization_id 
		FROM organization_responsible org 
		JOIN bids b ON org.organization_id = b.author_id
		JOIN employee e ON e.id = org.user_id
		WHERE e.username = $1 AND b.id = $2`

	err = p.db.QueryRowContext(ctx, userCheckQuery, params.Username, bidId).Scan(&orgId)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotAllowed
	} else if err != nil {
		return nil, err
	}

	// 3. Обновляем статус предложения
	updateQuery := `
		UPDATE bids SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 RETURNING id, name, description, status, tender_id, author_type, author_id, created_at, updated_at`

	row := p.db.QueryRowContext(ctx, updateQuery, params.Status, bidId)

	var bid models.Bid
	err = row.Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Если статус "Approved", то закрываем тендер
	if params.Status == repos.BidStatus("Approved") {
		closeTenderQuery := `UPDATE tenders SET status = 'Closed', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
		_, err = p.db.ExecContext(ctx, closeTenderQuery, bid.TenderId)
		if err != nil {
			return nil, err
		}
	}

	return &bid, nil
}

func (p *Postgres) EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (*models.Bid, error) {
	var bid models.Bid
	var authorId repos.BidAuthorId

	// Проверяем существование бида
	err := p.db.QueryRowContext(ctx, `SELECT id, author_id, name, description, status, tender_id, author_type, created_at FROM bids WHERE id = $1`, bidId).Scan(
		&bid.Id, &bid.AuthorId, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bid not found")
		}
		return nil, err
	}

	// Получаем id пользователя по username
	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, username).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Проверяем, что автор бида совпадает с пользователем
	if bid.AuthorId != authorId {
		return nil, fmt.Errorf("you are not the author of this bid")
	}

	// 3. Вносим изменения
	if params.Name != nil {
		bid.Name = *params.Name
	}
	if params.Description != nil {
		bid.Description = *params.Description
	}

	_, err = p.db.ExecContext(ctx, `
        UPDATE bids 
        SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3`, bid.Name, bid.Description, bidId)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновленный bid
	return &bid, nil
}

func (p *Postgres) GetBidsByUsername(ctx context.Context, username repos.Username) ([]*models.Bid, error) {
	var bids []*models.Bid
	var authorId int

	// Получаем id пользователя по username
	err := p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, username).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Получаем все биды с этим author_id
	rows, err := p.db.QueryContext(ctx, `SELECT id, name, description, status, tender_id, author_type, created_at FROM bids WHERE author_id = $1`, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Собираем результаты
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (p *Postgres) GetBidByID(ctx context.Context, bidId repos.BidId) (*models.Bid, error) {
	var bid models.Bid

	// Получаем bid по ID
	err := p.db.QueryRowContext(ctx, `SELECT id, name, description, status, tender_id, author_type, author_id, created_at FROM bids WHERE id = $1`, bidId).Scan(
		&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.AuthorId, &bid.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bid not found")
		}
		return nil, err
	}

	return &bid, nil
}

func (p *Postgres) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (*models.Bid, error) {
	var bid models.Bid
	var authorId repos.BidAuthorId

	// 1. Проверяем существование бида
	err := p.db.QueryRowContext(ctx, `SELECT id, author_id, status, tender_id FROM bids WHERE id = $1`, bidId).Scan(&bid.Id, &bid.AuthorId, &bid.Status, &bid.TenderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bid not found")
		}
		return nil, err
	}

	// 2. Проверяем, является ли пользователь автором бида
	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.Username).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if bid.AuthorId != authorId {
		return nil, fmt.Errorf("you are not the author of this bid")
	}

	// 3. Меняем статус бида
	_, err = p.db.ExecContext(ctx, `UPDATE bids SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, params.Decision, bidId)
	if err != nil {
		return nil, err
	}

	// Обновляем данные в модели
	bid.Status = models.BidStatus(params.Decision)

	return &bid, nil
}

func (p *Postgres) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) ([]*models.BidReview, error) {
	var tenderExists bool
	var authorId, responsibleId repos.OrganizationId
	var reviews []*models.BidReview

	// 1. Валидируем существование тендера
	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tenders WHERE id = $1)`, tenderId).Scan(&tenderExists)
	if err != nil {
		return nil, err
	}
	if !tenderExists {
		return nil, fmt.Errorf("tender not found")
	}

	// 2. Проверяем, существует ли author username и есть ли у него биды
	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.AuthorUsername).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}

	var hasBids bool
	err = p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE author_id = $1 AND tender_id = $2)`, authorId, tenderId).Scan(&hasBids)
	if err != nil {
		return nil, err
	}
	if !hasBids {
		return nil, fmt.Errorf("no bids found for the author")
	}

	// 3. Валидируем пользователя как ответственного за организацию и владельца тендера
	err = p.db.QueryRowContext(ctx, `
        SELECT o.id 
        FROM organization_responsible r
        JOIN tenders t ON t.organization_id = r.organization_id
        JOIN employee e ON e.id = r.user_id
        WHERE t.id = $1 AND e.username = $2`, tenderId, params.RequesterUsername).Scan(&responsibleId)
	if err != nil {
		return nil, fmt.Errorf("you are not authorized to view reviews for this tender")
	}

	// 4. Получаем список отзывов на биды
	rows, err := p.db.QueryContext(ctx, `
        SELECT f.id, f.description, f.created_at 
        FROM bid_feedbacks f 
        JOIN bids b ON b.id = f.bid_id 
        WHERE b.author_id = $1 AND b.tender_id = $2`, authorId, tenderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feedback models.BidReview
		err := rows.Scan(&feedback.Id, &feedback.Description, &feedback.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &feedback)
	}

	return reviews, nil
}

func (p *Postgres) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (*models.Bid, error) {
	var bidExists bool
	var authorId int
	var bid models.Bid

	// 1. Валидируем существование бида
	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE id = $1)`, bidId).Scan(&bidExists)
	if err != nil {
		return nil, err
	}
	if !bidExists {
		return nil, fmt.Errorf("bid not found")
	}

	// 2. Проверяем существование юзера в базе
	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.Username).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 3. Оставляем отзыв
	_, err = p.db.ExecContext(ctx, `
        INSERT INTO bid_feedbacks (bid_id, author_id, description, created_at) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`, bidId, authorId, params.BidFeedback)
	if err != nil {
		return nil, err
	}

	// Возвращаем обновленный bid
	err = p.db.QueryRowContext(ctx, `SELECT id, name, description, status, tender_id, author_id, created_at FROM bids WHERE id = $1`, bidId).Scan(
		&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorId, &bid.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (p *Postgres) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (*models.Bid, error) {
	var bidExists bool
	var currentVersion int32
	var bid models.Bid

	// 1. Валидируем существование бида
	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE id = $1)`, bidId).Scan(&bidExists)
	if err != nil {
		return nil, err
	}
	if !bidExists {
		return nil, fmt.Errorf("bid not found")
	}

	// 2. Проверяем существование версии для этого бида
	err = p.db.QueryRowContext(ctx, `
        SELECT version_number, name, description, status 
        FROM bid_versions 
        WHERE bid_id = $1 AND version_number = $2`, bidId, version).Scan(&currentVersion, &bid.Name, &bid.Description, &bid.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version not found")
		}
		return nil, err
	}

	// 3. Откатываемся к этой версии
	_, err = p.db.ExecContext(ctx, `
        UPDATE bids 
        SET name = $1, description = $2, status = $3, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $4`, bid.Name, bid.Description, bid.Status, bidId)
	if err != nil {
		return nil, err
	}

	// 4. Фиксируем новую запись как текущую версию
	_, err = p.db.ExecContext(ctx, `
        INSERT INTO bid_versions (bid_id, version_number, name, description, status, created_at, is_current) 
        VALUES ($1, (SELECT COALESCE(MAX(version_number), 0) + 1 FROM bid_versions WHERE bid_id = $1), $2, $3, $4, CURRENT_TIMESTAMP, TRUE)`,
		bidId, bid.Name, bid.Description, bid.Status)
	if err != nil {
		return nil, err
	}

	return &bid, nil
}
