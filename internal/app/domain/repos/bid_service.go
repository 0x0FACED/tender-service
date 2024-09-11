package repos

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/models"
)

// Username Уникальный slug пользователя.
type Username = string

// PaginationLimit defines model for paginationLimit.
type PaginationLimit = int32

// PaginationOffset defines model for paginationOffset.
type PaginationOffset = int32

// BidAuthorId Уникальный идентификатор автора предложения, присвоенный сервером.
type BidAuthorId = string

// BidAuthorType Тип автора
type BidAuthorType string

// BidDecision Решение по предложению
type BidDecision string

// BidDescription Описание предложения
type BidDescription = string

// BidFeedback Отзыв на предложение
type BidFeedback = string

// BidId Уникальный идентификатор предложения, присвоенный сервером.
type BidId = string

// BidName Полное название предложения
type BidName = string

// BidReviewDescription Описание предложения
type BidReviewDescription = string

// BidReviewId Уникальный идентификатор отзыва, присвоенный сервером.
type BidReviewId = string

// BidStatus Статус предложения
type BidStatus string

// BidVersion Номер версии посел правок
type BidVersion = int32

// BidService предоставляет методы для работы с предложениями.
type BidService interface {
	CreateBid(ctx context.Context, params CreateBidParams) (models.Bid, error)
	GetUserBids(ctx context.Context, params GetUserBidsParams) ([]*models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderId TenderId, params GetBidsForTenderParams) ([]*models.Bid, error)
	GetBidStatus(ctx context.Context, bidId BidId, params GetBidStatusParams) (BidStatus, error)
	UpdateBidStatus(ctx context.Context, bidId BidId, params UpdateBidStatusParams) (models.Bid, error)
	EditBid(ctx context.Context, bidId BidId, username Username, params EditBidParams) (models.Bid, error)
	SubmitBidDecision(ctx context.Context, bidId BidId, params SubmitBidDecisionParams) (models.Bid, error)
	SubmitBidFeedback(ctx context.Context, bidId BidId, params SubmitBidFeedbackParams) (models.Bid, error)
	RollbackBid(ctx context.Context, bidId BidId, version int32, params RollbackBidParams) (models.Bid, error)
	GetBidReviews(ctx context.Context, tenderId TenderId, params GetBidReviewsParams) ([]*models.BidReview, error)
}

// CreateBidParams определяет параметры для создания нового предложения.
type CreateBidParams struct {
	Name            *BidName        `json:"name"`
	Description     *BidDescription `json:"description"`
	Status          *BidStatus      `json:"status"`
	TenderID        *TenderId       `json:"tenderId"`
	OrganizationID  *OrganizationId `json:"organizationId"`
	CreatorUsername *Username       `json:"creatorUsername"`
}

// BidReview Отзыв о предложении
type BidReview struct {
	// CreatedAt Серверная дата и время в момент, когда пользователь отправил отзыв на предложение.
	// Передается в формате RFC3339.
	CreatedAt string `json:"createdAt"`

	// Description Описание предложения
	Description BidReviewDescription `json:"description"`

	// Id Уникальный идентификатор отзыва, присвоенный сервером.
	Id BidReviewId `json:"id"`
}

// GetUserBidsParams defines parameters for GetUserBids.
type GetUserBidsParams struct {
	// Limit Максимальное число возвращаемых объектов. Используется для запросов с пагинацией.
	//
	// Сервер должен возвращать максимальное допустимое число объектов.
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Какое количество объектов должно быть пропущено с начала. Используется для запросов с пагинацией.
	Offset   *PaginationOffset `form:"offset,omitempty" json:"offset,omitempty"`
	Username *Username         `form:"username,omitempty" json:"username,omitempty"`
}

// EditBidParams defines parameters for EditBid.
type EditBidParams struct {
	Name        *BidName        `json:"name,omitempty"`
	Description *BidDescription `json:"description,omitempty"`
}

// SubmitBidFeedbackParams defines parameters for SubmitBidFeedback.
type SubmitBidFeedbackParams struct {
	BidFeedback BidFeedback `form:"bidFeedback" json:"bidFeedback"`
	Username    Username    `form:"username" json:"username"`
}

// RollbackBidParams defines parameters for RollbackBid.
type RollbackBidParams struct {
	Username Username `form:"username" json:"username"`
}

// GetBidStatusParams defines parameters for GetBidStatus.
type GetBidStatusParams struct {
	Username Username `form:"username" json:"username"`
}

// UpdateBidStatusParams defines parameters for UpdateBidStatus.
type UpdateBidStatusParams struct {
	Status   BidStatus `form:"status" json:"status"`
	Username Username  `form:"username" json:"username"`
}

// SubmitBidDecisionParams defines parameters for SubmitBidDecision.
type SubmitBidDecisionParams struct {
	Decision BidDecision `form:"decision" json:"decision"`
	Username Username    `form:"username" json:"username"`
}

// GetBidsForTenderParams defines parameters for GetBidsForTender.
type GetBidsForTenderParams struct {
	Username Username `form:"username" json:"username"`

	// Limit Максимальное число возвращаемых объектов. Используется для запросов с пагинацией.
	//
	// Сервер должен возвращать максимальное допустимое число объектов.
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Какое количество объектов должно быть пропущено с начала. Используется для запросов с пагинацией.
	Offset *PaginationOffset `form:"offset,omitempty" json:"offset,omitempty"`
}

// GetBidReviewsParams defines parameters for GetBidReviews.
type GetBidReviewsParams struct {
	// AuthorUsername Имя пользователя автора предложений, отзывы на которые нужно просмотреть.
	AuthorUsername Username `form:"authorUsername" json:"authorUsername"`

	// RequesterUsername Имя пользователя, который запрашивает отзывы.
	RequesterUsername Username `form:"requesterUsername" json:"requesterUsername"`

	// Limit Максимальное число возвращаемых объектов. Используется для запросов с пагинацией.
	//
	// Сервер должен возвращать максимальное допустимое число объектов.
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Какое количество объектов должно быть пропущено с начала. Используется для запросов с пагинацией.
	Offset *PaginationOffset `form:"offset,omitempty" json:"offset,omitempty"`
}
