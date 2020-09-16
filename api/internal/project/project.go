package project

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

// The Project package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound         = errors.New("project not found")
	ErrInvalidID        = errors.New("id provided was not a valid UUID")
	ErrEmptyColumnOrder = errors.New("project column order provided was empty")
)

func Retrieve(ctx context.Context, repo *database.Repository, pId string, uId string) (*Project, error) {
	var p Project

	if _, err := uuid.Parse(pId); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"project_id",
		"name",
		"open",
		"user_id",
		"column_order",
		"created",
	).From(
		"projects",
	).Where(sq.Eq{"project_id": "?", "user_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &p, q, pId, uId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func List(ctx context.Context, repo *database.Repository, uId string) ([]Project, error) {
	var p Project
	var ps = make([]Project, 0)

	stmt := repo.SQ.Select(
		"project_id",
		"name",
		"open",
		"user_id",
		"column_order",
		"created",
	).From("projects").Where(sq.Eq{"user_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	rows, err := repo.DB.QueryContext(ctx, q, uId)
	if err != nil {
		return nil, errors.Wrap(err, "selecting projects")
	}
	for rows.Next() {
		err = rows.Scan(&p, &p.Name, &p.Open, &p.UserID, (*pq.StringArray)(&p.ColumnOrder), &p.Created)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row into Struct")
		}
		ps = append(ps, p)
	}

	return ps, nil
}

// Create adds a new Project
func Create(ctx context.Context, repo *database.Repository, np NewProject, uId string, now time.Time) (*Project, error) {
	p := Project{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Open:        true,
		UserID:      uId,
		ColumnOrder: []string{"column-1", "column-2", "column-3", "column-4"},
		Created:     now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"projects",
	).SetMap(map[string]interface{}{
		"project_id":   p.ID,
		"name":         p.Name,
		"open":         p.Open,
		"user_id":      p.UserID,
		"column_order": pq.Array(p.ColumnOrder),
		"created":      now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting project: %v", np)
	}

	return &p, nil
}

// Update modifies data about a Project. It will error if the specified ID is
// invalid or does not reference an existing Project.
func Update(ctx context.Context, repo *database.Repository, pId string, update UpdateProject, uId string) error {
	p, err := Retrieve(ctx, repo, pId, uId)
	if err != nil {
		return err
	}

	p.Name = update.Name
	p.Open = update.Open

	if len(update.ColumnOrder) != 0 {
		return ErrEmptyColumnOrder
	}

	p.ColumnOrder = update.ColumnOrder

	stmt := repo.SQ.Update(
		"project",
	).SetMap(map[string]interface{}{
		"name":         p.Name,
		"open":         p.Open,
		"column_order": p.ColumnOrder,
	}).Where(sq.Eq{"project_id": pId,"user_id": uId})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating project")
	}

	return nil
}

// Delete removes the Project identified by a given ID.
func Delete(ctx context.Context, repo *database.Repository, pId, uId string) error {
	if _, err := uuid.Parse(pId); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"projects",
	).Where(sq.Eq{"project_id": pId, "user_id": uId})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting project %s", pId)
	}

	return nil
}
