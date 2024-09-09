package server

import "github.com/0x0FACED/tender-service/internal/app/domain/repos"

// Псевдоним типа для CreateBid запроса
type CreateBidJSONRequestBody CreateBidJSONBody

// Псевдоним типа для EditBid запроса
type EditBidJSONRequestBody EditBidJSONBody

// Псевдоним типа для CreateTender запроса
type CreateTenderJSONRequestBody CreateTenderJSONBody

// Псевдоним типа для EditTender запроса
type EditTenderJSONRequestBody EditTenderJSONBody

// CreateBidJSONBody defines parameters for CreateBid.
type CreateBidJSONBody struct {
	// CreatorUsername Уникальный slug пользователя.
	CreatorUsername repos.Username `json:"creatorUsername"`

	// Description Описание предложения
	Description repos.BidDescription `json:"description"`

	// Name Полное название предложения
	Name repos.BidName `json:"name"`

	// OrganizationId Уникальный идентификатор организации, присвоенный сервером.
	OrganizationId repos.OrganizationId `json:"organizationId"`

	// Status Статус предложения
	Status repos.BidStatus `json:"status"`

	// TenderId Уникальный идентификатор тендера, присвоенный сервером.
	TenderId repos.TenderId `json:"tenderId"`
}

type EditBidJSONBody struct {
	Name        *repos.BidName        `json:"name,omitempty"`
	Description *repos.BidDescription `json:"description,omitempty"`
}

// CreateTenderJSONBody defines parameters for CreateTender.
type CreateTenderJSONBody struct {
	// CreatorUsername Уникальный slug пользователя.
	CreatorUsername repos.Username `json:"creatorUsername"`

	// Description Описание тендера
	Description repos.TenderDescription `json:"description"`

	// Name Полное название тендера
	Name repos.TenderName `json:"name"`

	// OrganizationId Уникальный идентификатор организации, присвоенный сервером.
	OrganizationId repos.OrganizationId `json:"organizationId"`

	// ServiceType Вид услуги, к которой относиться тендер
	ServiceType repos.TenderServiceType `json:"serviceType"`

	// Status Статус тендер
	Status repos.TenderStatus `json:"status"`
}

// EditTenderJSONBody defines parameters for EditTender.
type EditTenderJSONBody struct {
	// Description Описание тендера
	Description *repos.TenderDescription `json:"description,omitempty"`

	// Name Полное название тендера
	Name *repos.TenderName `json:"name,omitempty"`

	// ServiceType Вид услуги, к которой относиться тендер
	ServiceType *repos.TenderServiceType `json:"serviceType,omitempty"`
}
