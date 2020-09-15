package column

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

// The Column package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound  = errors.New("column not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo *database.Repository, id string) (*Column, error) {
	var c Column

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"column_id",
		"project_id",
		"title",
		"column",
		"task_ids",
		"created",
	).From(
		"columns",
	).Where(sq.Eq{"column_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &c, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func List(ctx context.Context, repo *database.Repository) ([]Column, error) {
	var c []Column

	stmt := repo.SQ.Select(
		"column_id",
		"project_id",
		"title",
		"column",
		"task_ids",
		"created",
	).From("columns").Where(sq.Eq{"project_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &c, q); err != nil {
		return nil, errors.Wrap(err, "selecting columns")
	}

	return c, nil
}

// Create adds a new Column
func Create(ctx context.Context, repo *database.Repository, nc NewColumn, now time.Time) (*Column, error) {

	c := Column{
		ID:        uuid.New().String(),
		ProjectID: nc.ProjectID,
		Title:     nc.Title,
		Column:    nc.Column,
		TaskIDS:   nc.TaskIDS,
		Created:   now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"columns",
	).SetMap(map[string]interface{}{
		"column_id":  c.ID,
		"project_id": c.ProjectID,
		"title":      c.Title,
		"column":     c.Column,
		"task_ids":   c.TaskIDS,
		"created":    now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting column: %v", nc)
	}

	return &c, nil
}

// Update modifies data about a Column. It will error if the specified ID is
// invalid or does not reference an existing Column.
func Update(ctx context.Context, repo *database.Repository, id string, update UpdateColumn) error {
	c, err := Retrieve(ctx, repo, id)
	if err != nil {
		return err
	}

	c.Title = update.Title
	c.TaskIDS = update.TaskIDS

	stmt := repo.SQ.Update(
		"columns",
	).SetMap(map[string]interface{}{
		"title":        c.Title,
		"task_ids":        c.TaskIDS,
	}).Where(sq.Eq{"column_id": id, "project_id": c.ProjectID})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating column")
	}

	return nil
}

// Delete removes the column identified by a given ID.
func Delete(ctx context.Context, repo *database.Repository, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"columns",
	).Where(sq.Eq{"column_id": id})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting column %s", id)
	}

	return nil
}
