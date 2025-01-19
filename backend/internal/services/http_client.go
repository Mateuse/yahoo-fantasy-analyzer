package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP error: %d - %s", e.StatusCode, e.Message)
}

// HttpRequest handles making HTTP requests
func HttpRequest(req *http.Request) (map[string]interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HttpError{StatusCode: resp.StatusCode, Message: string(body)}
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

func HttpXMLRequest(req *http.Request) (map[string]interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	// Convert XML to JSON
	jsonResponse, err := utils.XMLtoJSON(body)
	if err != nil {
		return nil, fmt.Errorf("error converting XML to JSON: %w", err)
	}
	cleanedResponse, err := utils.RemoveFantasyContent(jsonResponse)
	if err != nil {
		return nil, err
	}

	return cleanedResponse, nil
}

func AuthHttpXMLRequest(sessionId, url string) (map[string]interface{}, error) {
	//Get access token
	accessToken, err := GetAuthToken(sessionId)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, utils.NewNotFoundError(fmt.Sprintf("No session found for user: %s", sessionId))
		}
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}

	// make get request with token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := HttpXMLRequest(req)
	if err == nil {
		return resp, nil // Success
	}

	// check for token expiry 401
	if respErr, ok := err.(*HttpError); ok && respErr.StatusCode == http.StatusUnauthorized {
		// Get user id for sql
		userId, err := GetUserId(sessionId)
		if err != nil {
			return nil, err
		}

		//reresh token
		accessToken, refreshErr := ExchangeRefreshToken(userId)
		if refreshErr != nil {
			if utils.IsNotFoundError(refreshErr) {
				return nil, utils.NewNotFoundError(fmt.Sprintf("No refresh token found for user: %s", userId))
			}
			return nil, fmt.Errorf("failed to get refresh token: %w", refreshErr)
		}

		// retry  GET request with the new access token
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err = HttpRequest(req)
		if err != nil {
			return nil, fmt.Errorf("request failed after token refresh: %w", err)
		}
		return resp, nil
	}

	return nil, fmt.Errorf("request failed: %w", err) //other errors
}
