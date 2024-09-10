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
	// валидируем все поля параметров
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
		// TODO: error handling
		return models.Bid{}, err
	}
	return *bid, nil
}

func (b *BidServiceImpl) GetUserBids(ctx context.Context, params repos.GetUserBidsParams) ([]*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) GetBidsForTender(ctx context.Context, tenderId repos.TenderId, params repos.GetBidsForTenderParams) ([]*models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) GetBidStatus(ctx context.Context, bidId repos.BidId, params repos.GetBidStatusParams) (repos.BidStatus, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) UpdateBidStatus(ctx context.Context, bidId repos.BidId, params repos.UpdateBidStatusParams) (models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) EditBid(ctx context.Context, bidId repos.BidId, username repos.Username, params repos.EditBidParams) (models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) SubmitBidDecision(ctx context.Context, bidId repos.BidId, params repos.SubmitBidDecisionParams) (models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) SubmitBidFeedback(ctx context.Context, bidId repos.BidId, params repos.SubmitBidFeedbackParams) (models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) RollbackBid(ctx context.Context, bidId repos.BidId, version int32, params repos.RollbackBidParams) (models.Bid, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BidServiceImpl) GetBidReviews(ctx context.Context, tenderId repos.TenderId, params repos.GetBidReviewsParams) (models.BidReview, error) {
	panic("not implemented") // TODO: Implement
}
