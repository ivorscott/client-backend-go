package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/project"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/pkg/errors"
)

// Project holds the application state needed by the handler methods.
type Projects struct {
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

// List gets all Project
func (p *Projects) List(w http.ResponseWriter, r *http.Request) error {
	id := p.auth0.GetUserById(r)

	list, err := project.List(r.Context(), p.repo, id)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, list, http.StatusOK)
}

// Retrieve a single Project
func (p *Projects) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pId := chi.URLParam(r, "id")
	uId := p.auth0.GetUserById(r)

	ts, err := project.Retrieve(r.Context(), p.repo, pId, uId)
	if err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for projects %q", pId)
		}
	}

	return web.Respond(r.Context(), w, ts, http.StatusOK)
}

// Create a new Project
func (p *Projects) Create(w http.ResponseWriter, r *http.Request) error {
	uId := p.auth0.GetUserById(r)

	var np project.NewProject
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	ts, err := project.Create(r.Context(), p.repo, np, uId, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, ts, http.StatusCreated)
}

// Update decodes the body of a request to update an existing project. The ID
// of the project is part of the request URL.
func (p *Projects) Update(w http.ResponseWriter, r *http.Request) error {
	pId := chi.URLParam(r, "id")
	uId := p.auth0.GetUserById(r)

	var update project.UpdateProject
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding project update")
	}

	if err := project.Update(r.Context(), p.repo, pId, update, uId); err != nil {
		switch err {
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pId)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}

// Delete removes a single Project identified by an ID in the request URL.
func (p *Projects) Delete(w http.ResponseWriter, r *http.Request) error {
	pId := chi.URLParam(r, "id")
	uId := p.auth0.GetUserById(r)

	if err := project.Delete(r.Context(), p.repo, pId, uId); err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting project %q", pId)
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusNoContent)
}
