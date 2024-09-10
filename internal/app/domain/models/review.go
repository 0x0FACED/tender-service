package models

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

// BidReviewDescription Описание предложения
type BidReviewDescription = string

// BidReviewId Уникальный идентификатор отзыва, присвоенный сервером.
type BidReviewId = string
