package servicesimpl

import (
	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
	e "github.com/0x0FACED/tender-service/internal/app/errs"
)

func validateCreateBid(params repos.CreateBidParams) *e.ServiceError {
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

func validateGetUserBids(params repos.GetUserBidsParams) *e.ServiceError {
	if *params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateGetBidsForTender(params repos.GetBidsForTenderParams) *e.ServiceError {
	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

// 'Created', 'Published', 'Canceled', 'Approved', 'Rejected'
func validateGetBidStatus(params repos.GetBidStatusParams) *e.ServiceError {
	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateUpdateBidStatus(params repos.UpdateBidStatusParams) *e.ServiceError {
	if params.Status == "Created" {
		err := e.New("already exists status", e.ErrAlreadyExists)
		return err
	}

	if params.Status != "Published" && params.Status != "Canceled" && params.Status != "Approved" && params.Status != "Rejected" {
		err := e.New("unknown status", e.ErrUnknownStatus)
		return err
	}

	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateEditBid(params repos.EditBidParams) *e.ServiceError {
	if len(*params.Name) > MAX_BID_NAME_SIZE {
		err := e.New("name length exceeded", e.ErrExceededLength)
		return err
	}

	if len(*params.Description) > MAX_BID_DESCRIPTION_SIZE {
		err := e.New("desc length exceeded", e.ErrExceededLength)
		return err
	}

	return nil
}

func validateSubmitBidDecision(params repos.SubmitBidDecisionParams) *e.ServiceError {
	if params.Decision != "Approved" && params.Decision != "Rejected" {
		err := e.New("unknown decision", e.ErrUnknownDecision)
		return err
	}

	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateSubmitBidFeedback(params repos.SubmitBidFeedbackParams) *e.ServiceError {
	if len(params.BidFeedback) > 1000 {
		err := e.New("feedback length exceeded", e.ErrExceededLength)
		return err
	}

	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateRollbackBid(params repos.RollbackBidParams) *e.ServiceError {
	if params.Username == "" {
		err := e.New("empty username", e.ErrEmpty)
		return err
	}

	return nil
}

func validateGetBidReviews(params repos.GetBidReviewsParams) *e.ServiceError {
	if params.AuthorUsername == "" {
		err := e.New("empty author username", e.ErrEmpty)
		return err
	}

	if params.RequesterUsername == "" {
		err := e.New("empty requester username", e.ErrEmpty)
		return err
	}

	return nil
}
