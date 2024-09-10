package servicesimpl

import (
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	e "github.com/0x0FACED/tender-service/internal/app/errs"
)

func ValidateCreateBid(params repos.CreateBidParams) *e.ServiceError {
	// Проверяем длину имена, если вдруг > 100
	if len(*params.Name) > MAX_BID_NAME_SIZE {
		err := e.New("name length exceeded", e.ErrExceededLength)
		return err
	}

	// Проверяем длина описания, если вдруг > 500
	if len(*params.Description) > MAX_BID_DESCRIPTION_SIZE {
		err := e.New("desc length exceeded", e.ErrExceededLength)
		return err
	}

	// При создании статус может быть только Created, что логично
	if *params.Status != "Created" {
		err := e.New("invalid status, must be Created", e.ErrInvalidStatusCreateBid)
		return err
	}

	// На этом этапе мы просто проверяем, пустой ли айди тендера
	if *params.TenderID == "" {
		err := e.New("tender id is empty", e.ErrEmpty)
		return err
	}

	// OrganizationId не проверяем, потому что биды могут создавать НЕ от имени организации, судя по спецификации
	// А значит он может быть пуст

	// Чекаем, пустой ли username
	if *params.CreatorUsername == "" {
		err := e.New("desc length exceeded", e.ErrExceededLength)
		return err
	}

	return nil
}
