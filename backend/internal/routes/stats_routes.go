package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterStatsRoutes(router *mux.Router) {
	router.HandleFunc("/get-fantasy-league-player-stats/league/{leagueId}/player/{playerId}", handlers.GetFantasyLeaguePlayerStats).Methods("GET")
}
