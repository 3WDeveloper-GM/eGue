package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/3WDeveloper-GM/eGue/cmd/config"
	"github.com/3WDeveloper-GM/eGue/cmd/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routesHandler(app *config.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/api/healthcheck", handlers.GetHealthCheckHandler(app))
	r.Post("/api/search", handlers.PostSearchEmailMatchphrase(app))
	return r
}

func startServer(app *config.Application) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		IdleTimeout:  time.Minute,
		Handler:      routesHandler(app),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.ListenAndServe()
}
