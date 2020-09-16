package user

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
)

// The User package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound       = errors.New("user not found")
	ErrInvalidID      = errors.New("id provided was not a valid UUID")
)

// Retrieve finds the User identified by a given Auth0ID.
func RetrieveMeById(ctx context.Context, repo *database.Repository, userId string) (*User, error) {
	var u User

	if _, err := uuid.Parse(userId); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"user_id",
		"auth0_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"locale",
		"picture",
		"created",
	).From(
		"users",
	).Where(sq.Eq{"user_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &u, q, userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

// Retrieve finds the User identified by a given Auth0ID.
func RetrieveMeBySubject(ctx context.Context, repo *database.Repository, auth0Id string) (*User, error) {
	var u User

	stmt := repo.SQ.Select(
		"user_id",
		"auth0_id",
		"email",
		"first_name",
		"last_name",
		"email_verified",
		"locale",
		"picture",
		"created",
	).From(
		"users",
	).Where(sq.Eq{"auth0_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &u, q, auth0Id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

// Create adds a new User
func Create(ctx context.Context, repo *database.Repository, nu NewUser, auth0Id string, now time.Time) (*User, error) {

	u := User{
		ID:            uuid.New().String(),
		Auth0ID:       auth0Id,
		Email:         nu.Email,
		EmailVerified: nu.EmailVerified,
		FirstName:     nu.FirstName,
		LastName:      nu.LastName,
		Picture:       nu.Picture,
		Locale:        nu.Locale,
		Created:       now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"user_id":        u.ID,
		"auth0_id":       u.Auth0ID,
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"first_name":     u.FirstName,
		"last_name":      u.LastName,
		"picture":        u.Picture,
		"locale":         u.Locale,
		"created":        u.Created,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting user: %v", nu)
	}

	return &u, nil
}
