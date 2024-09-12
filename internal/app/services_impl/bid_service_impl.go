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

const (
	MAX_BID_NAME_SIZE        = 100
	MAX_BID_DESCRIPTION_SIZE = 500
)

// Всякие валидации здесь будут и вызовы БД
type BidServiceImpl struct {
	db database.BidRepository
}

func NewBidService(db database.BidRepository) repos.BidService {
	return &BidServiceImpl{
		db: db,
	}
}

func (b *BidServiceImpl) CreateBid(ctx context.Context, params repos.CreateBidParams) (models.Bid, error) {
	if err := validateCreateBid(params); err != nil {
		return models.Bid{}, err.Error()
	}

	bid, err := b.db.CreateBid(ctx, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error) {

	if err := validateGetUserBids(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetUserBids(ctx, params)
}

func (b *BidServiceImpl) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error) {

	if err := validateGetBidsForTender(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetBidsForTender(ctx, tenderId, params)
}

func (b *BidServiceImpl) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error) {

	if err := validateGetBidStatus(params); err != nil {
		return "", err.Error()
	}
	return b.db.GetBidStatus(ctx, bidId, params)
}

func (b *BidServiceImpl) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (models.Bid, error) {

	if err := validateUpdateBidStatus(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.UpdateBidStatus(ctx, bidId, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (models.Bid, error) {

	if err := validateEditBid(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.EditBid(ctx, bidId, username, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (models.Bid, error) {

	if err := validateSubmitBidDecision(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.SubmitBidDecision(ctx, bidId, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (models.Bid, error) {
	if err := validateSubmitBidFeedback(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.SubmitBidFeedback(ctx, bidId, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (models.Bid, error) {
	if err := validateRollbackBid(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.RollbackBid(ctx, bidId, version, params)
	if err != nil {
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) ([]*models.BidReview, error) {
	if err := validateGetBidReviews(params); err != nil {
		return nil, err.Error()
	}
	reviews, err := b.db.GetBidReviews(ctx, tenderId, params)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
