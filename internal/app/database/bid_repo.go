package database

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

type BidRepository interface {
	CreateBid(ctx context.Context, params repos.CreateBidParams) (*models.Bid, error)
	GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error)
	GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error)
	UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (*models.Bid, error)
	EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (*models.Bid, error)
	GetBidsByUsername(ctx context.Context, username repos.Username) ([]*models.Bid, error)
	GetBidByID(ctx context.Context, bidId repos.BidId) (*models.Bid, error)

	BidDecisionRepository
	BidFeedbackRepository
	BidVersionRepository
}

type BidDecisionRepository interface {
	SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (*models.Bid, error)
}

type BidFeedbackRepository interface {
	// Review == feedback
	GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) (*models.BidReview, error)
	SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (*models.Bid, error)
}

type BidVersionRepository interface {
	RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (*models.Bid, error)
}
