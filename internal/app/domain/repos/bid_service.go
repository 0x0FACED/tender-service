package repos

import "github.com/labstack/echo/v4"

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
	// Получение списка ваших предложений
	GetUserBids(ctx echo.Context, params GetUserBidsParams) error
	// Создание нового предложения
	CreateBid(ctx echo.Context, params CreateBidParams) error
	// Редактирование параметров предложения
	EditBid(ctx echo.Context, bidId BidId, params EditBidParams) error
	// Отправка отзыва по предложению
	SubmitBidFeedback(ctx echo.Context, bidId BidId, params SubmitBidFeedbackParams) error
	// Откат версии предложения
	RollbackBid(ctx echo.Context, bidId BidId, version int32, params RollbackBidParams) error
	// Получение текущего статуса предложения
	GetBidStatus(ctx echo.Context, bidId BidId, params GetBidStatusParams) error
	// Изменение статуса предложения
	UpdateBidStatus(ctx echo.Context, bidId BidId, params UpdateBidStatusParams) error
	// Отправка решения по предложению
	SubmitBidDecision(ctx echo.Context, bidId BidId, params SubmitBidDecisionParams) error
	GetBidsForTender(ctx echo.Context, tenderId TenderId, params GetBidsForTenderParams) error
	// Просмотр отзывов на прошлые предложения
	// (GET /bids/{tenderId}/reviews)
	GetBidReviews(ctx echo.Context, tenderId TenderId, params GetBidReviewsParams) error
}

// CreateBidParams определяет параметры для создания нового предложения.
type CreateBidParams struct {
	Name            BidName        `json:"name"`
	Description     BidDescription `json:"description"`
	Status          BidStatus      `json:"status"`
	TenderID        TenderId       `json:"tenderId"`
	OrganizationID  OrganizationId `json:"organizationId"`
	CreatorUsername Username       `json:"creatorUsername"`
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
	Username Username `form:"username" json:"username"`
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
