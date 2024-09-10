package postgres

import (
	"context"
	"database/sql"
	"errors"

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
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) GetBidsByUsername(ctx context.Context, username repos.Username) ([]*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) GetBidByID(ctx context.Context, bidId repos.BidId) (*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) (*models.BidReview, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (p *Postgres) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}
