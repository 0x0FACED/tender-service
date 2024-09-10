package e

import "errors"

// Ошикби на уровне валидации входных данных

type ServiceError struct {
	err error
	msg string
}

func New(msg string, err error) *ServiceError {
	return &ServiceError{
		err: err,
		msg: msg,
	}
}

func (s *ServiceError) Error() error {
	return s.err
}

func (s *ServiceError) Message() string {
	return s.msg
}

var (
	ErrExceededLength         = errors.New("exceeded length")
	ErrEmpty                  = errors.New("empty")
	ErrInvalidStatusCreateBid = errors.New("must be Created only")
)
