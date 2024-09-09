package servicesimpl

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
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

func (b *BidServiceImpl) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) CreateBid(ctx context.Context, params repos.CreateBidParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (models.Bid, error) {
	panic("impl me")
}

func (b *BidServiceImpl) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) error {
	panic("impl me")
}
