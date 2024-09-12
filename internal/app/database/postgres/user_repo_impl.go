package postgres

import (
	"context"
	"database/sql"

	"github.com/0x0FACED/tender-service/internal/app/domain/repos"
)

func (p *Postgres) GetUserIDByUsername(ctx context.Context, username repos.Username) (int, error) {
	var userID int
	userQuery := `SELECT id FROM employee WHERE username = $1`
	err := p.db.QueryRowContext(ctx, userQuery, username).Scan(&userID)
	if err == sql.ErrNoRows {
		return -1, ErrUserNotFound
	} else if err != nil {
		return -1, err
	}

	return userID, nil
}
