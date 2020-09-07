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

// The product package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

// Retrieve finds the product identified by a given ID.
func RetrieveMe(ctx context.Context, repo *database.Repository, id string) (*User, error) {
	var u User

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"id",
		"name",
		"price",
		"description",
		"created",
		"tags",
	).From(
		"users",
	).Where(sq.Eq{"id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

// Create adds a new Product
func Create(ctx context.Context, repo *database.Repository, nu NewUser, now time.Time) (*User, error) {
	u := User{
		ID:          uuid.New().String(),
		Auth0ID:        nu.Auth0ID,
		Email:       nu.Email,
		Picture: nu.Picture,
		Created:     now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"users",
	).SetMap(map[string]interface{}{
		"id":      u.ID,
		"auth0ID": u.Auth0ID,
		"email":   u.Email,
		"picture": u.Picture,
		"created": u.Created,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}

	return &u, nil
}
