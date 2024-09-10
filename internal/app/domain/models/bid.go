package models

// Bid Информация о предложении
type Bid struct {
	// AuthorId Уникальный идентификатор автора предложения, присвоенный сервером.
	AuthorId BidAuthorId `json:"authorId"`

	// AuthorType Тип автора
	AuthorType BidAuthorType `json:"authorType"`

	// CreatedAt Серверная дата и время в момент, когда пользователь отправил предложение на создание.
	// Передается в формате RFC3339.
	CreatedAt string `json:"createdAt"`

	// Description Описание предложения
	Description BidDescription `json:"description"`

	// Id Уникальный идентификатор предложения, присвоенный сервером.
	Id BidId `json:"id"`

	// Name Полное название предложения
	Name BidName `json:"name"`

	// Status Статус предложения
	Status BidStatus `json:"status"`

	// TenderId Уникальный идентификатор тендера, присвоенный сервером.
	TenderId TenderId `json:"tenderId"`

	// Version Номер версии посел правок
	Version BidVersion `json:"version"`
}

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

// BidStatus Статус предложения
type BidStatus string

// BidVersion Номер версии посел правок
type BidVersion = int32
