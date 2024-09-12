package database

import (
	"context"

	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

type UserRepository interface {
	GetUserIDByUsername(ctx context.Context, username repos.Username) (int, error)
}
