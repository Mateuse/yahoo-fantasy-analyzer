package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterNHLRoutes(router *mux.Router) {
	router.HandleFunc("/save-all-teams-schedule", handlers.SaveAllTeamsSchedule).Methods("GET")
	router.HandleFunc("/get-next-game/{team}", handlers.GetTeamNextGameDate).Methods("GET")
	router.HandleFunc("/get-player-game-stats/{playerId}", handlers.GetPlayerGameStats).Methods("GET")
	router.HandleFunc("/get-player-game-stats/{playerId}/season/{season}", handlers.GetPlayerGameStats).Methods("GET")
	router.HandleFunc("/get-team-roster/{teamAbrev}", handlers.GetTeamRoster).Methods("GET")
	router.HandleFunc("/get-team-roster/{teamAbrev}/{season}", handlers.GetTeamRoster).Methods("GET")
	router.HandleFunc("/map-players", handlers.SavePlayerIDMapping).Methods("POST")
}
