package servicesimpl

import (
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	e "github.com/0x0FACED/tender-service/internal/app/errs"
)

var (
	MAX_TENDER_NAME_SIZE        = 100
	MAX_TENDER_DESCRIPTION_SIZE = 1000
)

func validateCreateTender(params repos.CreateTenderParams) *e.ServiceError {
	if params.Name == nil {
		err := e.New("empty tender name", e.ErrEmpty)
		return err
	}

	if len(*params.Name) > MAX_TENDER_NAME_SIZE {
		err := e.New("name length exceeded", e.ErrExceededLength)
		return err
	}

	if params.Description == nil {
		err := e.New("empty tender desc", e.ErrEmpty)
		return err
	}

	if len(*params.Description) > MAX_TENDER_DESCRIPTION_SIZE {
		err := e.New("desc length exceeded", e.ErrExceededLength)
		return err
	}

	if *params.Status != "Created" && *params.Status != "Published" && *params.Status != "Closed" {
		err := e.New("unknown tender status", e.ErrUnknownStatus)
		return err
	}

	if *params.ServiceType != "Construction" && *params.ServiceType != "Delivery" && *params.ServiceType != "Manufacture" {
		err := e.New("unknown tender status", e.ErrUnknownStatus)
		return err
	}

	if params.OrganizationID == nil {
		err := e.New("empty org id", e.ErrEmpty)
		return err
	}

	if params.CreatorUsername == nil {
		err := e.New("empty creator username", e.ErrEmpty)
		return err
	}

	return nil
}

// status = 'Created', 'Published', 'Closed'
// type = 'Construction', 'Delivery', 'Manufacture'
func validateGetTenders(params repos.GetTendersParams) *e.ServiceError {
	// TODO: валидация слайса типов
	return nil
}
func validateGetUserTenders(params repos.GetUserTendersParams) *e.ServiceError {
	if params.Username == nil {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}
	return nil
}
func validateEditTender(params repos.EditTenderParams) *e.ServiceError {
	if params.Name == nil {
		err := e.New("empty tender name", e.ErrEmpty)
		return err
	}

	if len(*params.Name) > MAX_TENDER_NAME_SIZE {
		err := e.New("name length exceeded", e.ErrExceededLength)
		return err
	}

	if params.Description == nil {
		err := e.New("empty tender desc", e.ErrEmpty)
		return err
	}

	if len(*params.Description) > MAX_TENDER_DESCRIPTION_SIZE {
		err := e.New("desc length exceeded", e.ErrExceededLength)
		return err
	}

	if params.ServiceType == nil {
		err := e.New("empty tender type", e.ErrEmpty)
		return err
	}

	if *params.ServiceType != "Construction" && *params.ServiceType != "Delivery" && *params.ServiceType != "Manufacture" {
		err := e.New("unknown tender status", e.ErrUnknownStatus)
		return err
	}
	return nil
}
func validateRollbackTender(params repos.RollbackTenderParams) *e.ServiceError {
	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}
	return nil
}
func validateGetTenderStatus(params repos.GetTenderStatusParams) *e.ServiceError {
	if *params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}
	return nil
}
func validateUpdateTenderStatus(params repos.UpdateTenderStatusParams) *e.ServiceError {
	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	if params.Status == "Created" && params.Status != "Published" && params.Status != "Closed" {
		err := e.New("unknown tender status", e.ErrUnknownStatus)
		return err
	}

	return nil
}
