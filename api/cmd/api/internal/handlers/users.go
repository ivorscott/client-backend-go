package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
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
	repo  *database.Repository
	log   *log.Logger
	auth0 *mid.Auth0
}

// Retrieve a single User
func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	sub := u.auth0.GetAuthTokenSubject(r)
	us, err := user.RetrieveMe(r.Context(), u.repo, sub)

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

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

// Create a new User
func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {
	// (Problem) Auth0's access token doesn't reference our internal user's id in the database,
	// causing the backend to make an extra request to retrieve the current's id in our system even
	// when the user's record is not needed for the request.
	//
	// (Solution)
	// Store the management API token in our internal database

	// 1. create a new token if it doesn't exists or expired, store it in database
	// 3. call Auth0's management API with required credentials and token
	// 4. update the user.app_metadata in Auth0
	// (Frequency) Once on user creation
	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	us, err := user.Create(r.Context(), u.repo, nu, time.Now())
	if err != nil {
		return err
	}

	//token, err := u.auth0.GetManagementToken()
	//if err != nil {
	//	return err
	//}

	// use management api endpoint to update user's app_metadata
	// store user_id in app_metadata

	return web.Respond(r.Context(), w, us, http.StatusCreated)
}
