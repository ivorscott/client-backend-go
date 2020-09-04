package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"github.com/rs/cors"
)

func API(shutdown chan os.Signal, repo *database.Repository, log *log.Logger, FrontendAddress, Auth0Audience, Auth0Domain string) http.Handler {
	auth := mid.Auth0{Audience: Auth0Audience, Domain: Auth0Domain}
	app := web.NewApp(shutdown, log, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{FrontendAddress},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})

	{
		c := HealthCheck{repo: repo}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	p := Products{repo: repo, log: log}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

	return c.Handler(auth.Authenticate(app))
}
