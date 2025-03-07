package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func GetFantasyLeaguePlayerStats(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	leagueId := vars["leagueId"]

	if leagueId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing League Id", nil)
		return
	}

	playerId := vars["playerId"]

	if playerId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing Player Id", nil)
		return
	}

	cachedPlayerStats, err := services.GetCachedResponse(playerId+leagueId, "getplayerstats")

	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedPlayerStats != nil {
		utils.CustomResponse(w, http.StatusAccepted, "Successfully retrieved league player stats from cache", cachedPlayerStats)
		return
	}

	leagueOptions, err := services.GetLeagueSettings(userSession, leagueId)
	if err != nil {
		utils.CustomResponse(w, http.StatusBadRequest, "Failed Getting league options", err)
	}

	player, err := services.GetPlayerStats(userSession, playerId)
	if err != nil {
		utils.CustomResponse(w, http.StatusBadRequest, "Failed Getting Player stats", err)
	}

	leaguePlayerStats, totalPoints, err := services.GetLeaguePlayerStats(leagueOptions.StatModifiers, *player)
	if err != nil {
		utils.CustomResponse(w, http.StatusBadRequest, "Failed Getting Player stats for league", err)
		return
	}

	// Structure the response to include total points
	response := map[string]interface{}{
		"player_stats": leaguePlayerStats,
		"total_points": totalPoints,
	}

	err = services.CacheResponse(playerId+leagueId, "getPlayerStats", response, utils.GetTTL())
	if err != nil {
		fmt.Errorf("%w", err)
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved player stats for league", response)
}

func GetProjectedVsActual(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	fTeamId := vars["fTeamId"]

	if fTeamId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing Fantasy Team Id", nil)
		return
	}

	cachedStats, err := services.GetCachedResponse(fTeamId, "getProjectedvsExpected")

	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedStats != nil {
		utils.CustomResponse(w, http.StatusOK, "Successfully retrieved Projected vs Expected Stats from cache", cachedStats)
		return
	}

}
