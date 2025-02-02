package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func YahooLogin(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	redirectURI := os.Getenv("YAHOO_REDIRECT_URI")

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Missing OAuth configuration", http.StatusInternalServerError)
		return
	}

	// Construct Yahoo OAuth URL
	authURL := fmt.Sprintf(
		"https://api.login.yahoo.com/oauth2/request_auth?client_id=%s&redirect_uri=%s&response_type=code&scope=fspt-r",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
	)

	// Redirect the user to Yahoo OAuth
	http.Redirect(w, r, authURL, http.StatusFound)
}

func YahooCallback(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization code from the query parameters
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	// Call ExchangeAuthCode to exchange the authorization code for tokens
	userid, err := services.ExchangeAuthCode(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange authorization code: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct the redirect URL with userid as a query parameter
	frontendURL := "http://localhost:5173"
	redirectURL, err := url.Parse(frontendURL)
	if err != nil {
		http.Error(w, "Invalid frontend URL", http.StatusInternalServerError)
		return
	}

	// Add userid as a query parameter
	query := redirectURL.Query()
	query.Set("userid", userid)
	redirectURL.RawQuery = query.Encode()

	// Redirect the user to the constructed URL
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

type ClearCacheRequest struct {
	Operation string `json:"operation"` // Optional: Specifies which cache operation to clear
}

func ClearCache(w http.ResponseWriter, r *http.Request) {
	var req ClearCacheRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := services.ClearCache(req.Operation)
	if err != nil {
		utils.CustomResponse(w, http.StatusInternalServerError, "Failed to clear cache", err.Error())
		return
	}

	message := "Cache cleared successfully"
	if req.Operation != "" {
		message = "Cache cleared succesfully for operation " + req.Operation
	}

	utils.CustomResponse(w, http.StatusOK, message, nil)
}
