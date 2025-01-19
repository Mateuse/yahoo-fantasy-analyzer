package services

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
)

func GetUserLeagues(userSession string) (map[string]interface{}, error) {

	url := "https://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games/leagues"

	leaguesResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, err
	}

	leagues, err := ExtractLeaguesFromResponse(leaguesResponse)
	if err != nil {
		return nil, err
	}

	return leagues, nil
}

func GetLeague(userSession, leagueId string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s", leagueId)

	leagueResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, err
	}
	return leagueResponse, nil
}

func GetLeagueSettings(userSession, leagueId string) (*models.League, error) {
	// Check the database for existing league settings
	dbSettings, err := GetLeagueSettingsFromDB(leagueId)
	if err != nil {
		return nil, fmt.Errorf("error fetching league settings from database: %w", err)
	}

	if dbSettings != nil {
		// Convert DB settings to League object
		cachedSettings, err := MapToLeague(dbSettings)
		if err != nil {
			return nil, fmt.Errorf("error converting DB settings to league object: %w", err)
		}
		return cachedSettings, nil
	}

	// Make API call if not in cache
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s/settings", leagueId)
	leagueSettingsResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, fmt.Errorf("error fetching league settings from API: %w", err)
	}

	// Convert API response to League object
	leagueSettings, err := MapToLeague(leagueSettingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error converting API response to league: %w", err)
	}

	// Save settings to the database for future use
	err = SaveLeagueSettingsToDB(leagueId, leagueSettingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error saving league settings to database: %w", err)
	}

	return leagueSettings, nil
}

func GetLeagueSetting(userSession, leagueId, setting string) (string, error) {
	leagueSettings, err := GetLeagueSettings(userSession, leagueId)
	if err != nil {
		return "", err
	}

	// Get the reflection value of the League struct
	leagueValue := reflect.ValueOf(leagueSettings)

	// Dereference the pointer if necessary
	if leagueValue.Kind() == reflect.Ptr {
		leagueValue = leagueValue.Elem()
	}

	// Check if it's a struct after dereferencing
	if leagueValue.Kind() != reflect.Struct {
		return "", errors.New("leagueSettings is not a struct")
	}

	// Get the field by name (case-sensitive)
	field := leagueValue.FieldByName(setting)
	if !field.IsValid() {
		return "", fmt.Errorf("setting %s not found in league settings", setting)
	}

	// Convert the field value to a string based on its kind
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int, reflect.Int64:
		return strconv.Itoa(int(field.Int())), nil
	default:
		return "", errors.New("unsupported field type")
	}
}

func GetTeamWeekStats(sessionId, teamId, week string) (*models.Team, error) {
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/team/%s/stats;type=week;week=%s", teamId, week)

	teamWeeklyResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, err
	}

	teamWeekXMLResponse, err := MapToTeamWeek(teamWeeklyResponse)
	if err != nil {
		return nil, err
	}

	return teamWeekXMLResponse, nil
}

// func GetPlayerStats(sessionId, playerId string) *models
