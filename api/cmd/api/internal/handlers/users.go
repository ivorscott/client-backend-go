package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/ivorscott/devpie-client-backend-go/internal/user"
	"github.com/pkg/errors"
)

// Products holds the application state needed by the handler methods.
type Users struct {
	repo *database.Repository
	log  *log.Logger
}


// Retrieve a single Product
func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := user.RetrieveMe(r.Context(), u.repo, id)
	if err != nil {
		switch err {
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for products %q", id)
		}
	}

	return web.Respond(r.Context(), w, prod, http.StatusOK)
}

// Create a new Product
func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	prod, err := user.Create(r.Context(), u.repo, nu, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, prod, http.StatusCreated)
}
