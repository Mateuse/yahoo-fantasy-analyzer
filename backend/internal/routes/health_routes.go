package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterHealthRoutes(router *mux.Router) {
	router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
}
