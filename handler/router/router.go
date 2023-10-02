package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	healthzHandler := handler.NewHealthzHandler()
	TODOHandler := handler.NewTODOHandler(service.NewTODOService(todoDB))

	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/todos", TODOHandler)
	return mux
}
