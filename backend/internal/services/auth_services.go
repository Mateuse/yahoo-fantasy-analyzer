package services

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func getYahooAuthDetails() (string, string, string) {
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	clientSecret := os.Getenv("YAHOO_CLIENT_SECRET")
	tokenURL := "https://api.login.yahoo.com/oauth2/get_token"
	return clientID, clientSecret, tokenURL
}

func ExchangeAuthCode(authCode string) (string, error) {
	clientID, clientSecret, tokenURL := getYahooAuthDetails()

	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  "oob",
		"code":          authCode,
		"grant_type":    "authorization_code",
	}

	// Encode the payload as application/x-www-form-urlencoded
	formData := bytes.NewBufferString("")
	for key, value := range payload {
		formData.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}
	formData.Truncate(formData.Len() - 1) // Remove the trailing "&"

	// Make the POST request
	req, err := http.NewRequest("POST", tokenURL, formData)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResponse, err := HttpRequest(req)
	if err != nil {
		return "", err
	}

	userId, err := GetYahooUserProfile(tokenResponse["access_token"].(string))
	if err != nil {
		return "Error getting user profile", err
	}

	// Create a user session with access token
	err = CreateUserSession(userId, tokenResponse["access_token"].(string), tokenResponse["expires_in"].(float64))
	if err != nil {
		return "Error creating the user session in redis", err
	}

	err = repositories.AddRefreshToken(userId, tokenResponse["refresh_token"].(string))
	if err != nil {
		return "Error adding refresh token to db", err
	}

	return userId, nil
}

func GetAuthToken(sessionId string) (string, error) {
	session, err := GetUserSession(sessionId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			accessToken, refreshErr := ExchangeRefreshToken(sessionId)
			if refreshErr != nil {
				if utils.IsNotFoundError(refreshErr) {
					return "", utils.NewNotFoundError(fmt.Sprintf("No token found for user: %s", sessionId))
				}
				return "", fmt.Errorf("failed to get refresh token: %w", refreshErr)
			}

			return accessToken, nil
		}
		return "Error retrieving token", err
	}

	return session.AccessToken, nil
}

func GetUserId(sessionId string) (string, error) {
	session, err := GetUserSession(sessionId)
	if err != nil {
		return "Error retrieving token", err
	}

	return session.UserId, nil
}

func ExchangeRefreshToken(userId string) (string, error) {
	clientID, clientSecret, tokenURL := getYahooAuthDetails()

	refreshToken, err := repositories.GetRefreshToken(userId)
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	// Encode the payload as application/x-www-form-urlencoded
	formData := bytes.NewBufferString("")
	for key, value := range payload {
		formData.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}
	formData.Truncate(formData.Len() - 1) // Remove the trailing "&"

	// Make the POST request
	req, err := http.NewRequest("POST", tokenURL, formData)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResponse, err := HttpRequest(req)
	if err != nil {
		return "", err
	}

	// Create a user session with access token
	err = CreateUserSession(userId, tokenResponse["access_token"].(string), tokenResponse["expires_in"].(float64))
	if err != nil {
		return "Error creating the user session in redis", err
	}

	return tokenResponse["access_token"].(string), nil
}
