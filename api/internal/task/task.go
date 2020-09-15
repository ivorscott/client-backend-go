package task

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

// The Task package shouldn't know anything about http
// While it may identify common know errors, how to respond is left to the handlers
var (
	ErrNotFound  = errors.New("task not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)

func Retrieve(ctx context.Context, repo *database.Repository, id string) (*Task, error) {
	var t Task

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	stmt := repo.SQ.Select(
		"task_id",
		"title",
		"content",
		"created",
	).From(
		"tasks",
	).Where(sq.Eq{"task_id": "?"})

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &t, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

func List(ctx context.Context, repo *database.Repository) ([]Task, error) {
	var t = make([]Task, 0)

	stmt := repo.SQ.Select(
		"task_id",
		"title",
		"content",
		"created",
	).From("tasks").Where(sq.Eq{"task_id": "?"})
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.SelectContext(ctx, &t, q); err != nil {
		return nil, errors.Wrap(err, "selecting tasks")
	}

	return t, nil
}

// Create adds a new Task
func Create(ctx context.Context, repo *database.Repository, nt NewTask, now time.Time) (*Task, error) {

	t := Task{
		ID:      uuid.New().String(),
		Title:   nt.Title,
		Content: nt.Content,
		Created: now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"tasks",
	).SetMap(map[string]interface{}{
		"task_id": t.ID,
		"title":   t.Title,
		"content": t.Content,
		"created": now.UTC(),
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return nil, errors.Wrapf(err, "inserting tasks: %v", nt)
	}

	return &t, nil
}

// Update modifies data about a Task. It will error if the specified ID is
// invalid or does not reference an existing Task.
func Update(ctx context.Context, repo *database.Repository, id string, update UpdateTask) error {
	t, err := Retrieve(ctx, repo, id)
	if err != nil {
		return err
	}

	t.Title = update.Title
	t.Content = update.Content

	stmt := repo.SQ.Update(
		"tasks",
	).SetMap(map[string]interface{}{
		"title":    t.Title,
		"content": t.Content,
	}).Where(sq.Eq{"task_id": id})

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "updating task")
	}

	return nil
}

// Delete removes the task identified by a given ID.
func Delete(ctx context.Context, repo *database.Repository, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	stmt := repo.SQ.Delete(
		"tasks",
	).Where(sq.Eq{"task_id": id})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting task %s", id)
	}

	return nil
}
