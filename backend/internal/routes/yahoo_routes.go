package routes

import (
	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/handlers"
)

func RegisterYahooRoutes(router *mux.Router) {
	router.HandleFunc("/get-user-leagues", handlers.GetUserLeaguesHandler).Methods("GET")
	router.HandleFunc("/get-league-info/{leagueId}", handlers.GetLeagueInfo).Methods("GET")
	router.HandleFunc("/get-league-settings/{leagueId}", handlers.GetLeagueSettings).Methods("GET")
	router.HandleFunc("/get-team-weekly/team/{teamId}", handlers.GetTeamWeeklyStats).Methods("GET")
	router.HandleFunc("/get-player-stats/player/{playerId}", handlers.GetPlayerStats).Methods("GET")
	router.HandleFunc("/get-player-rank/league/{leagueId}/player/{playerId}", handlers.GetPlayerRankLeague).Methods("GET")
	router.HandleFunc("/get-all-players", handlers.GetAllPlayersYahoo).Methods("GET")
}
