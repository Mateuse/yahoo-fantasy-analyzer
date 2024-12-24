package services

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/responses"
)

func GetYahooUserProfile(accessToken string) (string, error) {
	url := "https://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user profile: status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the user ID from the response (simplified for demonstration)
	userID, nil := responses.ParseFantasyContent(body)
	if err != nil {
		return "", fmt.Errorf("failed to parse user profile: %w", err)
	}

	return userID, nil
}
