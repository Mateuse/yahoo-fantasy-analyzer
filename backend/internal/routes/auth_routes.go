package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", handlers.YahooLogin).Methods("GET")
	router.HandleFunc("/yahoo-redirect", handlers.YahooCallback).Methods("GET")
}
