package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ExchangeAuthCode(authCode string) (string, error) {
	clientID := os.Getenv("YAHOO_CLIENT_ID")
	clientSecret := os.Getenv("YAHOO_CLIENT_SECRET")

	tokenURL := "https://api.login.yahoo.com/oauth2/get_token"

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

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response: %s", string(body))
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	userId, err := GetYahooUserProfile(tokenResponse["access_token"].(string))
	if err != nil {
		return "Error getting user profile", err
	}

	// Create a user session with access token
	err = CreateUserSession(userId, tokenResponse["access_token"].(string), tokenResponse["expires_in"].(int))
	if err != nil {
		return "Error creating the user session in redis", err
	}

	return userId, nil
}

func SaveTokenToFile(filename string, tokenResponse map[string]interface{}) error {
	// Open or create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the token response as JSON and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Format with indentation
	if err := encoder.Encode(tokenResponse); err != nil {
		return err
	}

	return nil
}
