package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/google/uuid"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var (
	ErrTenderNotFound       = errors.New("tender not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrBidNotFound          = errors.New("bid not found")
	ErrVersionNotFound      = errors.New("version not found")
	ErrUserNotAllowed       = errors.New("user not allowed to view or update bid status")
	ErrNoBidsForAuthor      = errors.New("no bids found for the author")
	ErrNotAuthor            = errors.New("not author of the bid")
)

func (p *Postgres) CreateBid(ctx context.Context, params repos.CreateBidParams) (*models.Bid, error) {
	// Проверяем, существует ли тендер
	exists, err := p.IsTenderExists(ctx, *params.TenderID)
	if err != nil {
		p.logger.Error("Error in IsTenderExists()", zap.Any("params", params))
		return nil, err
	}

	if !exists {
		p.logger.Error("Tender not found")
		return nil, ErrTenderNotFound
	}

	// Если указан OrganizationID, проверяем существование организации
	if params.OrganizationID != nil {
		var orgExists bool
		orgQuery := `SELECT EXISTS (SELECT 1 FROM organization WHERE id = $1)`
		err := p.db.QueryRowContext(ctx, orgQuery, *params.OrganizationID).Scan(&orgExists)
		if err != nil {
			p.logger.Error("Error in check organization extists")
			return nil, err
		}

		if !orgExists {
			p.logger.Error("Ogranization not found")
			return nil, ErrOrganizationNotFound
		}
	}

	// Начинаем транзакцию, чтобы сохранить данные в обе таблицы атомарно
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("Error starting tx", zap.Error(err))
		return nil, err
	}

	// Откат транзакции при возникновении ошибки
	defer func() {
		if err != nil {
			p.logger.Error("Do Rollback tx", zap.Error(err))
			tx.Rollback()
		}
	}()

	// Создаем новое предложение (bid)
	bidQuery := `
		INSERT INTO bids (name, description, status, tender_id, author_type, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 
			(SELECT id FROM employee WHERE username = $6), 
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, description, status, tender_id, author_type, author_id, created_at`

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
		p.logger.Error("Error scanning vals to bid", zap.Error(err))
		return nil, err
	}

	versionQuery := `
		INSERT INTO bid_versions (bid_id, version_number, author_id, status, created_at, is_current)
		VALUES ($1, $2, (SELECT id FROM employee WHERE username = $3), $4, CURRENT_TIMESTAMP, TRUE)
		RETURNING version_number`

	err = tx.QueryRowContext(ctx, versionQuery,
		bid.Id,
		1,
		*params.CreatorUsername,
		*params.Status,
	).Scan(&bid.Version)
	if err != nil {
		p.logger.Error("Error query add version in tx", zap.Error(err))
		return nil, err
	}

	// Если все прошло успешно, фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		p.logger.Error("Error commit tx", zap.Error(err))
		return nil, err
	}

	return bid, nil
}

func (p *Postgres) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error) {
	// Проверяем, существует ли пользователь с указанным username
	userID, err := p.GetUserIDByUsername(ctx, *params.Username)
	if err != nil {
		p.logger.Error("Error GetUserIDByUsername()", zap.Error(err))
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
		p.logger.Error("Error get list of user bids", zap.Error(err))
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
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		bids = append(bids, &bid)
	}

	// Проверяем ошибки после итераций
	if err := rows.Err(); err != nil {
		p.logger.Error("Error rows.Err()", zap.Error(err))
		return nil, err
	}

	return bids, nil
}

func (p *Postgres) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error) {
	// Проверяем, существует ли тендер с указанным tenderId
	exists, err := p.IsTenderExists(ctx, tenderId)
	if err != nil {
		p.logger.Error("Error IsTenderExists()", zap.Error(err))
		return nil, ErrTenderNotFound
	}
	if !exists {
		p.logger.Error("Tender not found", zap.Any("params", params))
		return nil, ErrTenderNotFound
	}

	// Проверяем, существует ли пользователь с указанным username
	userID, err := p.GetUserIDByUsername(ctx, params.Username)
	if err != nil {
		p.logger.Error("Error GetUserIDByUsername()", zap.Error(err))
		return nil, ErrUserNotFound
	}

	// Формируем запрос для получения списка предложений для данного тендера
	bidQuery := `
		SELECT b.id, b.name, b.description, b.status, b.author_type, b.author_id, b.created_at, v.version_number
		FROM bids b
		JOIN bid_versions v ON b.id = v.bid_id
		JOIN tenders t ON b.tender_id = t.id
		WHERE b.tender_id = $1
		AND (
			b.author_id = $2
			OR
			t.organization_id = (SELECT organization_id FROM organization_responsible WHERE user_id = $2)
		)
		AND v.is_current = TRUE
		LIMIT $3 OFFSET $4
	`

	rows, err := p.db.QueryContext(ctx, bidQuery, tenderId, userID, params.Limit, params.Offset)
	if err != nil {
		p.logger.Error("Error get list of bids for tender()", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var bids []*models.Bid

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
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		bids = append(bids, &bid)
	}

	if err := rows.Err(); err != nil {
		p.logger.Error("Error rows.Err()", zap.Error(err))
		return nil, err
	}

	return bids, nil
}

func (p *Postgres) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error) {
	// Проверяем, существует ли предложение с данным bidId
	var status repos.BidStatus
	query := `SELECT status FROM bids WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, bidId).Scan(&status)
	if err == sql.ErrNoRows {
		p.logger.Error("Bid not found")
		return "", ErrBidNotFound
	} else if err != nil {
		p.logger.Error("Error get bid by id", zap.Error(err))
		return "", err
	}

	return status, nil
}

func (p *Postgres) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (*models.Bid, error) {
	var currentStatus repos.BidStatus
	query := `SELECT status FROM bids WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, bidId).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		p.logger.Error("Bid not found")
		return nil, ErrBidNotFound
	} else if err != nil {
		p.logger.Error("Error get bid by id", zap.Error(err))
		return nil, err
	}

	var orgId uuid.UUID
	userCheckQuery := `
		SELECT org.organization_id 
		FROM organization_responsible org 
		JOIN tenders t ON org.organization_id = t.organization_id
		JOIN bids b ON t.id = b.tender_id
		JOIN employee e ON e.id = org.user_id
		WHERE e.username = $1 AND b.id = $2`

	err = p.db.QueryRowContext(ctx, userCheckQuery, params.Username, bidId).Scan(&orgId)
	if err == sql.ErrNoRows {
		p.logger.Error("User not allowed")
		return nil, ErrUserNotAllowed
	} else if err != nil {
		p.logger.Error("Error user exists and check org responsible", zap.Error(err))
		return nil, err
	}

	updateQuery := `
		UPDATE bids SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 RETURNING id, name, description, status, tender_id, author_type, author_id, created_at`

	row := p.db.QueryRowContext(ctx, updateQuery, params.Status, bidId)

	var bid models.Bid
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
		p.logger.Error("Error row.Scan()", zap.Error(err))
		return nil, err
	}

	// Если статус "Approved", то закрываем тендер
	if params.Status == repos.BidStatus("Approved") {
		closeTenderQuery := `UPDATE tenders SET status = 'Closed', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
		_, err = p.db.ExecContext(ctx, closeTenderQuery, bid.TenderId)
		if err != nil {
			p.logger.Error("Error update bid to Closed status", zap.Error(err))
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
			p.logger.Error("Bid not found")
			return nil, ErrBidNotFound
		}
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, username).Scan(&authorId)
	if err != nil {
		p.logger.Error("User not found")
		return nil, ErrUserNotFound
	}

	if bid.AuthorId != authorId {
		p.logger.Error("Not author, declined")
		return nil, ErrNotAuthor
	}

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
		p.logger.Error("Error updating bid", zap.Error(err))
		return nil, err
	}

	var versionNumber int32
	err = p.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version_number), 0) + 1 FROM bid_versions WHERE bid_id = $1`, bidId).Scan(&versionNumber)
	if err != nil {
		p.logger.Error("Error fetching bid version", zap.Error(err))
		return nil, err
	}

	_, err = p.db.ExecContext(ctx, `
        INSERT INTO bid_versions (bid_id, version_number, author_id, status, created_at, is_current)
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, TRUE)`,
		bidId, versionNumber, bid.AuthorId, bid.Status)
	if err != nil {
		p.logger.Error("Error inserting bid version", zap.Error(err))
		return nil, err
	}

	bid.Version = versionNumber
	return &bid, nil
}

func (p *Postgres) GetBidsByUsername(ctx context.Context, username repos.Username) ([]*models.Bid, error) {
	var bids []*models.Bid
	var authorId int

	err := p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, username).Scan(&authorId)
	if err != nil {
		p.logger.Error("User not found")
		return nil, ErrUserNotFound
	}

	rows, err := p.db.QueryContext(ctx, `SELECT id, name, description, status, tender_id, author_type, created_at FROM bids WHERE author_id = $1`, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.CreatedAt)
		if err != nil {
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (p *Postgres) GetBidByID(ctx context.Context, bidId repos.BidId) (*models.Bid, error) {
	var bid models.Bid

	err := p.db.QueryRowContext(ctx, `SELECT id, name, description, status, tender_id, author_type, author_id, created_at FROM bids WHERE id = $1`, bidId).Scan(
		&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.AuthorId, &bid.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("Bid not found")
			return nil, ErrBidNotFound
		}
		return nil, err
	}

	return &bid, nil
}

func (p *Postgres) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (*models.Bid, error) {
	var bid models.Bid

	err := p.db.QueryRowContext(ctx, `SELECT id, author_id, status, tender_id FROM bids WHERE id = $1`, bidId).Scan(&bid.Id, &bid.AuthorId, &bid.Status, &bid.TenderId)
	if err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("Bid not found")
			return nil, ErrBidNotFound
		}
		return nil, err
	}

	var organizationId repos.OrganizationId
	err = p.db.QueryRowContext(ctx, `SELECT organization_id FROM tenders WHERE id = $1`, bid.TenderId).Scan(&organizationId)
	if err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("Tender not found")
			return nil, ErrTenderNotFound
		}
		return nil, err
	}

	var userId repos.BidAuthorId
	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.Username).Scan(&userId)
	if err != nil {
		p.logger.Error("User not found")
		return nil, ErrUserNotFound
	}

	var orgId repos.OrganizationId
	userCheckQuery := `
		SELECT o.organization_id
		FROM organization_responsible o
		JOIN employee e ON o.user_id = e.id
		WHERE e.id = $1 AND o.organization_id = $2`

	err = p.db.QueryRowContext(ctx, userCheckQuery, userId, organizationId).Scan(&orgId)
	if err == sql.ErrNoRows {
		p.logger.Error("User not authorized to make a decision")
		return nil, ErrUserNotAllowed
	} else if err != nil {
		p.logger.Error("Error checking user authorization", zap.Error(err))
		return nil, err
	}

	updateQuery := `
		UPDATE bids 
		SET status = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2 
		RETURNING id, author_id, name, description, status, tender_id, author_type, created_at`

	row := p.db.QueryRowContext(ctx, updateQuery, params.Decision, bidId)

	err = row.Scan(
		&bid.Id,
		&bid.AuthorId,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.CreatedAt,
	)
	if err != nil {
		p.logger.Error("Error updating bid status", zap.Error(err))
		return nil, err
	}

	var version int32
	insertVersionQuery := `
		INSERT INTO bid_versions (bid_id, version_number, status) 
		VALUES ($1, (SELECT COALESCE(MAX(version_number), 0) + 1 FROM bid_versions WHERE bid_id = $1), $2)
		RETURNING version_number`

	err = p.db.QueryRowContext(ctx, insertVersionQuery, bidId, params.Decision).Scan(&version)
	if err != nil {
		p.logger.Error("Error inserting bid version", zap.Error(err))
		return nil, err
	}

	bid.Version = version

	return &bid, nil
}

func (p *Postgres) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) ([]*models.BidReview, error) {
	var tenderExists bool
	var authorId, responsibleId repos.OrganizationId
	var reviews []*models.BidReview

	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tenders WHERE id = $1)`, tenderId).Scan(&tenderExists)
	if err != nil {
		p.logger.Error("Error check tender exists", zap.Error(err))
		return nil, err
	}
	if !tenderExists {
		p.logger.Error("Tender not found")
		return nil, ErrTenderNotFound
	}

	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.AuthorUsername).Scan(&authorId)
	if err != nil {
		p.logger.Error("User not found")
		return nil, ErrUserNotFound
	}

	var hasBids bool
	err = p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE author_id = $1 AND tender_id = $2)`, authorId, tenderId).Scan(&hasBids)
	if err != nil {
		p.logger.Error("Error check if bids exists", zap.Error(err))
		return nil, err
	}
	if !hasBids {
		p.logger.Error("No bids for author")
		return nil, ErrNoBidsForAuthor
	}

	// Валидируем пользователя как ответственного за организацию и владельца тендера
	err = p.db.QueryRowContext(ctx, `
        SELECT o.id 
        FROM organization_responsible r
        JOIN tenders t ON t.organization_id = r.organization_id
        JOIN employee e ON e.id = r.user_id
        WHERE t.id = $1 AND e.username = $2`, tenderId, params.RequesterUsername).Scan(&responsibleId)
	if err != nil {
		p.logger.Error("User not allowed")
		return nil, ErrUserNotAllowed
	}

	// Получаем список отзывов на биды
	rows, err := p.db.QueryContext(ctx, `
        SELECT f.id, f.description, f.created_at 
        FROM bid_feedbacks f 
        JOIN bids b ON b.id = f.bid_id 
        WHERE b.author_id = $1 AND b.tender_id = $2`, authorId, tenderId)
	if err != nil {
		p.logger.Error("Error get list of bid reviews", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var feedback models.BidReview
		err := rows.Scan(&feedback.Id, &feedback.Description, &feedback.CreatedAt)
		if err != nil {
			p.logger.Error("Error rows.Scan()", zap.Error(err))
			return nil, err
		}
		reviews = append(reviews, &feedback)
	}

	return reviews, nil
}

func (p *Postgres) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (*models.Bid, error) {
	var bidExists bool
	var authorId repos.BidAuthorId
	var bid models.Bid
	var versionNumber int32

	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE id = $1)`, bidId).Scan(&bidExists)
	if err != nil {
		p.logger.Error("Error check if bid exists", zap.Error(err))
		return nil, err
	}
	if !bidExists {
		p.logger.Error("Bid not found")
		return nil, ErrBidNotFound
	}

	err = p.db.QueryRowContext(ctx, `SELECT id FROM employee WHERE username = $1`, params.Username).Scan(&authorId)
	if err != nil {
		p.logger.Error("User not found")
		return nil, ErrUserNotFound
	}

	_, err = p.db.ExecContext(ctx, `
        INSERT INTO bid_feedbacks (bid_id, author_id, description, created_at) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`, bidId, authorId, params.BidFeedback)
	if err != nil {
		p.logger.Error("Error insert new review", zap.Error(err))
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, `SELECT version_number FROM bid_versions WHERE bid_id = $1 AND is_current = TRUE`, bidId).Scan(&versionNumber)
	if err != nil {
		p.logger.Error("Error retrieving current bid version", zap.Error(err))
		return nil, err
	}

	_, err = p.db.ExecContext(ctx, `UPDATE bids SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, bidId)
	if err != nil {
		p.logger.Error("Error update bid status", zap.Error(err))
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, `
		SELECT id, name, description, status, tender_id, author_id, author_type, created_at
		FROM bids WHERE id = $1`, bidId).Scan(
		&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorId, &bid.AuthorType, &bid.CreatedAt)
	if err != nil {
		p.logger.Error("Error get updated bid", zap.Error(err))
		return nil, err
	}

	bid.Version = versionNumber

	return &bid, nil
}

func (p *Postgres) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (*models.Bid, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var bidExists bool
	var bid models.Bid

	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM bids WHERE id = $1)`, bidId).Scan(&bidExists)
	if err != nil {
		p.logger.Error("Error check if bid exists", zap.Error(err))
		return nil, err
	}
	if !bidExists {
		p.logger.Error("Bid not found")
		return nil, ErrBidNotFound
	}

	err = tx.QueryRowContext(ctx, `
        SELECT status 
        FROM bid_versions 
        WHERE bid_id = $1 AND version_number = $2`, bidId, version).Scan(&bid.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			p.logger.Error("Version not found")
			return nil, ErrVersionNotFound
		}
		p.logger.Error("Error check if version exists", zap.Error(err))
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE bids 
        SET name = $1, description = $2, status = $3, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $4`, bid.Name, bid.Description, bid.Status, bidId)
	if err != nil {
		p.logger.Error("Error rollback to version", zap.Error(err))
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE bid_versions
        SET is_current = FALSE
        WHERE bid_id = $1 AND is_current = TRUE`, bidId)
	if err != nil {
		p.logger.Error("Error setting previous version to not current", zap.Error(err))
		return nil, err
	}

	// Создаем новую версию как текущую
	_, err = tx.ExecContext(ctx, `
        INSERT INTO bid_versions (bid_id, version_number, author_id, status, created_at, is_current) 
        VALUES ($1, 
                (SELECT COALESCE(MAX(version_number), 0) + 1 FROM bid_versions WHERE bid_id = $1), 
                (SELECT author_id FROM bids WHERE id = $1), 
                $2, CURRENT_TIMESTAMP, TRUE)`,
		bidId, bid.Status)
	if err != nil {
		p.logger.Error("Error creating new current version", zap.Error(err))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &bid, nil
}
