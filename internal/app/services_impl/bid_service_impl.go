package servicesimpl

import (
	"github.com/0x0FACED/tender-service/internal/app/database"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	"github.com/labstack/echo/v4"
)

type BidServiceImpl struct {
	db database.BidRepository
}

func NewBidService(db database.BidRepository) repos.BidService {
	return &BidServiceImpl{
		db: db,
	}
}

func (b *BidServiceImpl) GetUserBids(ctx echo.Context, params repos.GetUserBidsParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) CreateBid(ctx echo.Context, params repos.CreateBidParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) EditBid(ctx echo.Context, bidId repos.BidId, params repos.EditBidParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) SubmitBidFeedback(ctx echo.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) RollbackBid(ctx echo.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidStatus(ctx echo.Context, bidId repos.BidId, params repos.GetBidStatusParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) UpdateBidStatus(ctx echo.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) SubmitBidDecision(ctx echo.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidsForTender(ctx echo.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) error {
	panic("impl me")
}

func (b *BidServiceImpl) GetBidReviews(ctx echo.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) error {
	panic("impl me")
}
