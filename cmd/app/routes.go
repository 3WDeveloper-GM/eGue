package app

func (a *Application) setRoutes() {
	a.server.Get("/api/healthcheck", a.dependencies.searchHandler.HealthCheckHandler)
	a.server.Post("/api/_search", a.dependencies.searchHandler.SearchMails)
	a.server.Get("/api/_index", a.dependencies.indexHandler.IndexMails)
}
