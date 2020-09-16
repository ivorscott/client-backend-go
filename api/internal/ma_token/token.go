package ma_token

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound = errors.New("token not found")
)

// Retrieve auth0 management token
func Retrieve(ctx context.Context, repo *database.Repository) (*Token, error) {
	var t Token

	stmt := repo.SQ.Select(
		"ma_token_id",
		"token",
		"created",
	).From(
		"ma_token",
	).Limit(1)

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &t, q); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

// Persist auth0 management token
func Persist(ctx context.Context, repo *database.Repository, nt *Token, now time.Time) error {
	t := Token{
		ID:          uuid.New().String(),
		AccessToken: nt.AccessToken,
		Created:     now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"ma_token",
	).SetMap(map[string]interface{}{
		"ma_token_id": t.ID,
		"token": t.AccessToken,
		"created": t.Created,
	})

	if _, err := stmt.ExecContext(ctx); err != nil{
		return errors.Wrapf(err, "inserting token: %v", nt)
	}

	return nil
}

// Delete auth0 management token
func Delete(ctx context.Context, repo *database.Repository) error {
	stmt := repo.SQ.Delete("ma_token")

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting previous token")
	}

	return nil
}
