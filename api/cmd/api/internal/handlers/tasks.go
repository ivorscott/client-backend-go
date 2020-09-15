package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/ivorscott/devpie-client-backend-go/internal/task"
	"github.com/pkg/errors"
)

// Tasks holds the application state needed by the handler methods.
type Tasks struct {
	repo *database.Repository
	log  *log.Logger
	auth0 *mid.Auth0
}

// List gets all task
func (t *Tasks) List(w http.ResponseWriter, r *http.Request) error {
	list, err := task.List(r.Context(), t.repo)
	if err != nil {
		return err
	}

	if list == nil {
		list = []task.Task{}
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Task
func (t *Tasks) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	ts, err := task.Retrieve(r.Context(), t.repo, id)
	if err != nil {
		switch err {
		case task.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for tasks %q", id)
		}
	}

	return web.Respond(r.Context(), w, ts, http.StatusOK)
}

// Create a new Task
func (t *Tasks) Create(w http.ResponseWriter, r *http.Request) error {

	var nt task.NewTask
	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	ts, err := task.Create(r.Context(), t.repo, nt, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, ts, http.StatusCreated)
}

// Update decodes the body of a request to update an existing task. The ID
// of the task is part of the request URL.
func (t *Tasks) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update task.UpdateTask
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding task update")
	}

	if err := task.Update(r.Context(), t.repo, id, update); err != nil {
		switch err {
		case task.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating task %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single task identified by an ID in the request URL.
func (t *Tasks) Delete(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := task.Delete(r.Context(), t.repo, id); err != nil {
		switch err {
		case task.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting task %q", id)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}
