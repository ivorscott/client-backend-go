package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/ivorscott/devpie-client-backend-go/internal/user"
)

// Users holds the application state needed by the handler methods.
type Users struct {
	repo *database.Repository
	log  *log.Logger
}


// Retrieve a single User
func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	sub := fmt.Sprintf("%v",claims["sub"])

	prod, err := user.RetrieveMe(r.Context(), u.repo, sub)

	if err != nil {
		switch err {
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for user %q", sub)
		}
	}

	return web.Respond(r.Context(), w, prod, http.StatusOK)
	return nil
}

// Create a new User
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
