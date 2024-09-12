package server

import (
	"net/http"

	p "github.com/0x0FACED/tender-service/internal/app/database/postgres"
	e "github.com/0x0FACED/tender-service/internal/app/errs"
)

// ErrorResponse Используется для возвращения ошибки пользователю
type ErrorResponse struct {
	// Reason Описание ошибки в свободной форме
	Reason string `json:"reason"`
}

// Функция для получения статуса и сообщения об ошибке
func getStatusByError(err error) (int, ErrorResponse) {
	// Ошибки, связанные с некорректными данными со стороны пользователя
	switch err {
	case e.ErrExceededLength:
		return http.StatusBadRequest, ErrorResponse{Reason: "Превышена допустимая длина."}

	case e.ErrEmpty:
		return http.StatusBadRequest, ErrorResponse{Reason: "Пустое значение не допускается."}

	case e.ErrInvalidStatusCreateBid:
		return http.StatusBadRequest, ErrorResponse{Reason: "Некорректный статус. Статус должен быть 'Created'."}

	case e.ErrUnknownStatus:
		return http.StatusBadRequest, ErrorResponse{Reason: "Неизвестный статус."}

	case e.ErrUnknownDecision:
		return http.StatusBadRequest, ErrorResponse{Reason: "Неизвестное решение."}

	case e.ErrAlreadyExists:
		return http.StatusBadRequest, ErrorResponse{Reason: "Запись уже существует."}

	// Ошибки уровня сервиса/БД
	case p.ErrTenderNotFound:
		return http.StatusNotFound, ErrorResponse{Reason: "Тендер не найден."}

	case p.ErrOrganizationNotFound:
		return http.StatusNotFound, ErrorResponse{Reason: "Организация не найдена."}

	case p.ErrUserNotFound:
		return http.StatusNotFound, ErrorResponse{Reason: "Пользователь не найден."}

	case p.ErrBidNotFound:
		return http.StatusNotFound, ErrorResponse{Reason: "Заявка не найдена."}

	case p.ErrVersionNotFound:
		return http.StatusNotFound, ErrorResponse{Reason: "Версия не найдена."}

	case p.ErrUserNotAllowed:
		return http.StatusForbidden, ErrorResponse{Reason: "Пользователь не имеет прав на выполнение данного действия."}

	case p.ErrNotAuthor:
		return http.StatusForbidden, ErrorResponse{Reason: "Пользователь не является автором заявки."}

	case p.ErrNoBidsForAuthor:
		return http.StatusNoContent, ErrorResponse{Reason: "Заявки для данного автора не найдены."}
	}

	// Если ошибка не распознана (скорее всего при выполнении запроса к бд, но их сложно классифицировать)
	// поэтому 500 status
	return http.StatusInternalServerError, ErrorResponse{Reason: "Внутренняя ошибка сервера."}
}
