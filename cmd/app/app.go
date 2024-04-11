package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/handlers"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/jsonio"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/logger"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/service"
	indexer "github.com/3WDeveloper-GM/pipeline/cmd/pkg/service/Indexer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

const (
	portnumber = 4040
)

type Application struct {
	server       *chi.Mux
	client       *http.Client
	dependencies *dependency
	config       *ZsConfigration
}

type dependency struct {
	searchHandler *handlers.SearchHandler
	indexHandler  *handlers.IndexHandler
}

func (a *Application) setupConfig() {
	cfg := NewZsConfiguration()
	a.config = cfg
}

func (a *Application) setClient() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 64
	transport.MaxConnsPerHost = 128
	transport.MaxIdleConnsPerHost = 128

	client := &http.Client{
		Transport: transport,
	}

	a.client = client
}

func (a *Application) setupDepends() {

	logfile, err := os.Create("logfile.log")
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(logfile)

	adapter := adapter.NewAdapter(a.client, a.config)
	searchService := service.NewSearchMapper(a.config, *adapter)
	indexService := indexer.NewMailIndexer(a.config, adapter, logger)

	io := jsonio.NewJSONIO()

	shandler := handlers.NewSearchHandler(searchService, io)
	indexHandler := handlers.NewIndexHandler(indexService)

	a.dependencies = &dependency{
		searchHandler: shandler,
		indexHandler:  indexHandler,
	}
}

func (a *Application) setServer() {
	a.server = chi.NewRouter()
	a.server.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	a.server.Use(middleware.Logger)
	a.server.Use(middleware.Recoverer)
	a.server.Use(render.SetContentType(render.ContentTypeJSON))

	a.setRoutes()
}

func NewApplication() *Application {
	a := &Application{}
	a.setupConfig()
	a.setClient()
	a.setupDepends()
	a.setServer()

	return a
}

func (a *Application) StartApp() {
	log.Printf("Starting on http://localhost:%d", portnumber)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", portnumber), a.server); err != nil {
		panic(err)
	}

}
