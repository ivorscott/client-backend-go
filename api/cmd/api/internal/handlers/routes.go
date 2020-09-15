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


func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, FrontendAddress, Auth0Audience, Auth0Domain string) http.Handler {
	auth0 := &mid.Auth0{Audience: Auth0Audience, Domain: Auth0Domain}
	app := web.NewApp(shutdown, log, mid.Logger(log), auth0.Authenticate(), mid.Errors(log), mid.Panics(log))

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{FrontendAddress},
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "Content-Type", "Strict-Transport-Security"},
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

	app.Handle(http.MethodGet, "/v1/columns", c.List)
	app.Handle(http.MethodPost, "/v1/columns", c.Create)
	app.Handle(http.MethodGet, "/v1/columns/{id}", c.Retrieve)
	app.Handle(http.MethodPut, "/v1/columns/{id}", c.Update)
	app.Handle(http.MethodDelete, "/v1/columns/{id}", c.Delete)

	app.Handle(http.MethodGet, "/v1/tasks", t.List)
	app.Handle(http.MethodPost, "/v1/tasks", t.Create)
	app.Handle(http.MethodGet, "/v1/tasks/{id}", t.Retrieve)
	app.Handle(http.MethodPut, "/v1/tasks/{id}", t.Update)
	app.Handle(http.MethodDelete, "/v1/tasks/{id}", t.Delete)

	app.Handle(http.MethodGet, "/v1/projects", p.List)
	app.Handle(http.MethodPost, "/v1/projects", p.Create)
	app.Handle(http.MethodGet, "/v1/projects/{id}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/projects/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/projects/{id}", p.Delete)

	return cor.Handler(app)
}
