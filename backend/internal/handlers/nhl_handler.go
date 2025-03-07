package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func SaveAllTeamsSchedule(w http.ResponseWriter, r *http.Request) {
	if err := services.SaveAllTeamsSchedule(); err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to save schedule", err)
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully saved schedule in DB", nil)
}

func GetTeamNextGameDate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	team := vars["team"]

	nextGame, err := repositories.GetTeamNextGameDate(team)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Error", err)
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Next game for team is", nextGame)
}

func GetPlayerGameStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["playerId"]
	season := vars["season"]

	if playerId == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing player id", nil)
		return
	}

	if season == "" {
		season = utils.GetCurrentNhlSeason()
	}

	playerGameStats, err := services.GetPlayerGameStatsNHL(playerId, season)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Error getting player game stats", err)
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved player game stats", playerGameStats)
}

func GetTeamRoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamAbrev := vars["teamAbrev"]
	season := vars["season"]

	if teamAbrev == "" {
		utils.CustomResponse(w, http.StatusBadRequest, "Missing Team abreviation", nil)
		return
	}

	if season == "" {
		season = utils.GetCurrentNhlSeason()
	}

	roster, err := services.GetTeamRoster(teamAbrev, season)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to get team roster", err.Error())
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved team roster", roster)
}

func SavePlayerIDMapping(w http.ResponseWriter, r *http.Request) {

	err := services.MapNhlPlayerToYahoo()
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to map nhl and yahoo player ids", err.Error())
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully mapped nhl and yahoo player ids", nil)
}
