package servicesimpl

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

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
	if err := ValidateCreateBid(params); err != nil {
		return models.Bid{}, err.Error()
	}

	// Если все гуд, то добавляем в базу
	// Примечательно, что создавать bid можно как от имени организации, так и нет.
	// Если от имени организации, то authorId будет айди организации
	// И authorType будет 'Organization'

	// дальнейшая логика:
	/*
	  1. Проверяем, есть ли тендер в бд (CheckIfTenderExistsByID)
	  2. Проверяем наличие организации в бд (если не пустая строка)
	  3. Проверяем наличие юзера в бд.
	  4. Создаем бид.
	*/
	bid, err := b.db.CreateBid(ctx, params)
	if err != nil {

		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error) {

	if err := ValidateGetUserBids(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetUserBids(ctx, params)
}

func (b *BidServiceImpl) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error) {

	if err := ValidateGetBidsForTender(params); err != nil {
		return nil, err.Error()
	}
	return b.db.GetBidsForTender(ctx, tenderId, params)
}

func (b *BidServiceImpl) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error) {

	if err := ValidateGetBidStatus(params); err != nil {
		return "", err.Error()
	}
	return b.db.GetBidStatus(ctx, bidId, params)
}

func (b *BidServiceImpl) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (models.Bid, error) {

	if err := ValidateUpdateBidStatus(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.UpdateBidStatus(ctx, bidId, params)
	return *bid, err
}

func (b *BidServiceImpl) EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (models.Bid, error) {

	if err := ValidateEditBid(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.EditBid(ctx, bidId, username, params)
	return *bid, err
}

func (b *BidServiceImpl) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (models.Bid, error) {

	if err := ValidateSubmitBidDecision(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.SubmitBidDecision(ctx, bidId, params)
	return *bid, err
}

func (b *BidServiceImpl) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (models.Bid, error) {
	if err := ValidateSubmitBidFeedback(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.SubmitBidFeedback(ctx, bidId, params)
	return *bid, err
}

func (b *BidServiceImpl) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (models.Bid, error) {
	if err := ValidateRollbackBid(params); err != nil {
		return models.Bid{}, err.Error()
	}
	bid, err := b.db.RollbackBid(ctx, bidId, version, params)
	return *bid, err
}

func (b *BidServiceImpl) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) ([]*models.BidReview, error) {
	if err := ValidateGetBidReviews(params); err != nil {
		return nil, err.Error()
	}
	reviews, err := b.db.GetBidReviews(ctx, tenderId, params)
	return reviews, err
}
