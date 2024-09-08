package repos

import "github.com/labstack/echo/v4"

// OrganizationId Уникальный идентификатор организации, присвоенный сервером.
type OrganizationId = string

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

// TenderService предоставляет методы для работы с тендерами.
type TenderService interface {
	// Получение списка предложений для тендера
	GetBidsForTender(ctx echo.Context, tenderId TenderId, params GetBidsForTenderParams) error
	// Просмотр отзывов на прошлые предложения
	GetBidReviews(ctx echo.Context, tenderId TenderId, params GetBidReviewsParams) error
	// Получение списка тендеров
	GetTenders(ctx echo.Context, params GetTendersParams) error
	// Получение тендеров пользователя
	GetUserTenders(ctx echo.Context, params GetUserTendersParams) error
	// Создание нового тендера
	CreateTender(ctx echo.Context, params CreateTenderParams) error
	// Редактирование тендера
	EditTender(ctx echo.Context, tenderId TenderId, params EditTenderParams) error
	// Откат версии тендера
	RollbackTender(ctx echo.Context, tenderId TenderId, version int32, params RollbackTenderParams) error
	// Получение текущего статуса тендера
	GetTenderStatus(ctx echo.Context, tenderId TenderId, params GetTenderStatusParams) error
	// Изменение статуса тендера
	UpdateTenderStatus(ctx echo.Context, tenderId TenderId, params UpdateTenderStatusParams) error
}

type CreateTenderParams struct {
	Name            TenderName        `json:"name"`
	Description     TenderDescription `json:"description"`
	ServiceType     TenderServiceType `json:"serviceType"`
	Status          TenderStatus      `json:"status"`
	OrganizationID  OrganizationId    `json:"organizationId"`
	CreatorUsername Username          `json:"creatorUsername"`
}

// GetTendersParams defines parameters for GetTenders.
type GetTendersParams struct {
	// Limit Максимальное число возвращаемых объектов. Используется для запросов с пагинацией.
	//
	// Сервер должен возвращать максимальное допустимое число объектов.
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Какое количество объектов должно быть пропущено с начала. Используется для запросов с пагинацией.
	Offset *PaginationOffset `form:"offset,omitempty" json:"offset,omitempty"`

	// ServiceType Возвращенные тендеры должны соответствовать указанным видам услуг.
	//
	// Если список пустой, фильтры не применяются.
	ServiceType *[]TenderServiceType `form:"service_type,omitempty" json:"service_type,omitempty"`
}

// GetUserTendersParams defines parameters for GetUserTenders.
type GetUserTendersParams struct {
	// Limit Максимальное число возвращаемых объектов. Используется для запросов с пагинацией.
	//
	// Сервер должен возвращать максимальное допустимое число объектов.
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Какое количество объектов должно быть пропущено с начала. Используется для запросов с пагинацией.
	Offset   *PaginationOffset `form:"offset,omitempty" json:"offset,omitempty"`
	Username *Username         `form:"username,omitempty" json:"username,omitempty"`
}

// EditTenderParams defines parameters for EditTender.
type EditTenderParams struct {
	Username Username `form:"username" json:"username"`
}

// RollbackTenderParams defines parameters for RollbackTender.
type RollbackTenderParams struct {
	Username Username `form:"username" json:"username"`
}

// GetTenderStatusParams defines parameters for GetTenderStatus.
type GetTenderStatusParams struct {
	Username *Username `form:"username,omitempty" json:"username,omitempty"`
}

// UpdateTenderStatusParams defines parameters for UpdateTenderStatus.
type UpdateTenderStatusParams struct {
	Status   TenderStatus `form:"status" json:"status"`
	Username Username     `form:"username" json:"username"`
}
