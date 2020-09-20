package handlers

import (
	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, FrontendAddress,
	Auth0Audience, Auth0Domain, Auth0M2MClient, Auth0M2MSecret, AuthMAPIAudience string) http.Handler {

	auth0 := &mid.Auth0{
		Audience:     Auth0Audience,
		Domain:       Auth0Domain,
		M2MClient:    Auth0M2MClient,
		M2MSecret:    Auth0M2MSecret,
		MAPIAudience: AuthMAPIAudience,
	}

	app := web.NewApp(shutdown, log, mid.Logger(log), auth0.Authenticate(), mid.Errors(log), mid.Panics(log))

	cor := cors.New(cors.Options{
		AllowedOrigins: []string{FrontendAddress},
		AllowedHeaders: []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
		AllowedMethods: []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch,
		},
		AllowCredentials: true,
	})

	h := HealthCheck{repo: repo}

	app.Handle(http.MethodGet, "/v1/health", h.Health)

	u := Users{repo: repo, log: log, auth0: auth0}
	c := Columns{repo: repo, log: log, auth0: auth0}
	t := Tasks{repo: repo, log: log, auth0: auth0}
	p := Projects{repo: repo, log: log, auth0: auth0}

	app.Handle(http.MethodPost, "/v1/users", u.Create)
	app.Handle(http.MethodGet, "/v1/users/me", u.RetrieveMe)
	app.Handle(http.MethodGet, "/v1/projects", p.List)
	app.Handle(http.MethodPost, "/v1/projects", p.Create)
	app.Handle(http.MethodGet, "/v1/projects/{pid}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/projects/{pid}", p.Update)
	app.Handle(http.MethodDelete, "/v1/projects/{pid}", p.Delete)
	app.Handle(http.MethodGet, "/v1/projects/{pid}/columns", c.List)
	app.Handle(http.MethodGet, "/v1/projects/{pid}/tasks", t.List)
	app.Handle(http.MethodPost, "/v1/projects/{pid}/columns/{cid}/tasks", t.Create)
	app.Handle(http.MethodPatch, "/v1/projects/{pid}/tasks/{tid}", t.Update)
	app.Handle(http.MethodDelete, "/v1/projects/{pid}/columns/{cid}/tasks/{tid}", t.Delete)

	return cor.Handler(app)
}
