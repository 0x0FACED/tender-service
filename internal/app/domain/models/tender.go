package models

// OrganizationId Уникальный идентификатор организации, присвоенный сервером.
type OrganizationId = string

// Tender Информация о тендере
type Tender struct {
	// CreatedAt Серверная дата и время в момент, когда пользователь отправил тендер на создание.
	// Передается в формате RFC3339.
	CreatedAt string `json:"createdAt"`

	// Description Описание тендера
	Description TenderDescription `json:"description"`

	// Id Уникальный идентификатор тендера, присвоенный сервером.
	Id TenderId `json:"id"`

	// Name Полное название тендера
	Name TenderName `json:"name"`

	// OrganizationId Уникальный идентификатор организации, присвоенный сервером.
	OrganizationId OrganizationId `json:"organizationId"`

	// ServiceType Вид услуги, к которой относиться тендер
	ServiceType TenderServiceType `json:"serviceType"`

	// Status Статус тендер
	Status TenderStatus `json:"status"`

	// Version Номер версии посел правок
	Version TenderVersion `json:"version"`
}

// TenderDescription Описание тендера
type TenderDescription = string

// TenderId Уникальный идентификатор тендера, присвоенный сервером.
type TenderId = string

// TenderName Полное название тендера
type TenderName = string

// TenderServiceType Вид услуги, к которой относиться тендер
type TenderServiceType string

// TenderStatus Статус тендер
type TenderStatus string

// TenderVersion Номер версии посел правок
type TenderVersion = int32

// Username Уникальный slug пользователя.
type Username = string
