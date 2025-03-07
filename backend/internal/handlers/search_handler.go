package handlers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func GetPlayerByName(w http.ResponseWriter, r *http.Request) {

	userSession := r.Header.Get("user-session")
	if userSession == "" {
		utils.CustomResponse(w, http.StatusUnauthorized, "Missing user session", nil)
		return
	}

	vars := mux.Vars(r)

	playerName, err := url.QueryUnescape(vars["playerName"])
	if err != nil {
		utils.CustomResponse(w, http.StatusBadRequest, "Failed to get player name from query", nil)
		return
	}

	player, err := services.GetPlayerByName(userSession, playerName)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to get player details by name", err)
		return
	}

	utils.CustomResponse(w, http.StatusOK, "Successfully retrieved player details", player)
}
