package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func GetUserLeaguesHandler(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	cachedLeagues, err := services.GetCachedResponse(userSession, "getleagues")
	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedLeagues != nil {
		utils.CustomResponse(w, http.StatusOK, "Successfully retrieved user leagues from cache", cachedLeagues)
		return
	}

	leagues, err := services.GetUserLeagues(userSession)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.CustomResponse(w, http.StatusNotFound, err.Error(), nil)
		} else {
			utils.CustomResponse(w, http.StatusInternalServerError, "Failed to retrieve user leagues", err.Error())
		}
		return
	}

	// Define the TTL as the remaining time until the end of the day
	now := time.Now()
	expiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := time.Until(expiry)

	err = services.CacheResponse(userSession, "getleagues", leagues, ttl)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved user leagues", leagues)
}

func GetLeagueInfo(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	leagueId := vars["leagueId"]
	if leagueId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing league Id", nil)
		return
	}

	cachedLeague, err := services.GetCachedResponse(leagueId, "getleague")

	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedLeague != nil {
		utils.CustomResponse(w, http.StatusOK, "Successfully retrieved user leagues from cache", cachedLeague)
		return
	}

	league, err := services.GetLeague(userSession, leagueId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.CustomResponse(w, http.StatusNotFound, err.Error(), nil)
		} else {
			utils.CustomResponse(w, http.StatusInternalServerError, "Failed to retrieve user leagues", err.Error())
		}
		return
	}

	// Define the TTL as the remaining time until the end of the day
	now := time.Now()
	expiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := time.Until(expiry)

	err = services.CacheResponse(leagueId, "getleague", league, ttl)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved user leagues", league)
}

func GetLeagueSettings(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	leagueId := vars["leagueId"]
	if leagueId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing league Id", nil)
		return
	}

	// cachedLeagueSettings, err := services.GetCachedResponse(leagueId, "getleaguesettings")

	// if err != nil {
	// 	fmt.Errorf("%w", err)
	// }

	// if cachedLeagueSettings != nil {
	// 	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved league settings from cache", cachedLeagueSettings)
	// 	return
	// }

	leagueSettings, err := services.GetLeagueSettings(userSession, leagueId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.CustomResponse(w, http.StatusNotFound, err.Error(), err)
		} else {
			utils.CustomResponse(w, http.StatusInternalServerError, "Failed to fetch league settings", err.Error())
		}
		return
	}

	leagueSettingsMap, err := utils.StructToMap(leagueSettings)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to convert league settings", err.Error())
		return
	}

	// Define the TTL as the remaining time until the end of the day
	now := time.Now()
	expiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := time.Until(expiry)

	err = services.CacheResponse(leagueId, "getleaguesettings", leagueSettingsMap, ttl)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved league settings", leagueSettingsMap)

}

func GetTeamWeeklyStats(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	teamId := vars["teamId"]

	if teamId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing team id", nil)
		return
	}

	cachedWeeklyStats, err := services.GetCachedResponse(teamId, "getweeklystats")

	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedWeeklyStats != nil {
		utils.CustomResponse(w, http.StatusAccepted, "Successfully retrieved weekly stats from cache", cachedWeeklyStats)
		return
	}

	weeklyStats, err := services.GetTeamWeeklyStats(userSession, teamId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.CustomResponse(w, http.StatusNotFound, err.Error(), nil)
		} else {
			utils.CustomResponse(w, http.StatusInternalServerError, "Failed to retrieve team weekly stats", err.Error())
		}
		return
	}
	convertedWeeklyStats := utils.ConvertWeeklyStatsToMap(weeklyStats)
	// Define the TTL as the remaining time until the end of the day
	now := time.Now()
	expiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := time.Until(expiry)

	err = services.CacheResponse(teamId, "getweeklystats", convertedWeeklyStats, ttl)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved team weekly stats", convertedWeeklyStats)
}

func GetPlayerStats(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)
	playerId := vars["playerId"]

	if playerId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing Player Id", nil)
	}

	cachedPlayerStats, err := services.GetCachedResponse(playerId, "getplayerstats")

	if err != nil {
		fmt.Errorf("%w", err)
	}

	if cachedPlayerStats != nil {
		utils.CustomResponse(w, http.StatusAccepted, "Successfully retrieved player stats from cache", cachedPlayerStats)
		return
	}

}
