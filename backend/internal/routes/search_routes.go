package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterSearchRoutes(router *mux.Router) {
	router.HandleFunc("/get-player-by-name/player/{playerName}", handlers.GetPlayerByName).Methods("GET")
}
