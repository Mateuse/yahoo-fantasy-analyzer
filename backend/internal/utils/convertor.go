package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	xml2json "github.com/basgys/goxml2json"
)

func XMLtoJSON(body []byte) (map[string]interface{}, error) {
	if len(body) == 0 {
		return nil, errors.New("empty XML input")
	}

	// Validate XML by attempting to unmarshal into a generic struct
	var validateStruct interface{}
	if err := xml.Unmarshal(body, &validateStruct); err != nil {
		return nil, fmt.Errorf("invalid XML: %w", err)
	}

	// Convert XML to JSON
	jsonBuffer, err := xml2json.Convert(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error converting XML to JSON: %w", err)
	}

	// Unmarshal JSON string into a map
	var jsonObject map[string]interface{}
	if err := json.Unmarshal(jsonBuffer.Bytes(), &jsonObject); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return jsonObject, nil
}

func TeamtoLeagueId(teamId string) (string, error) {
	parts := strings.Split(teamId, ".t")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid team ID format: %s", teamId)
	}
	return parts[0], nil
}

func RemoveFantasyContent(response map[string]interface{}) (map[string]interface{}, error) {
	// Check if "fantasy_content" exists
	cleanedResponse, ok := response["fantasy_content"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'fantasy_content'")
	}

	return cleanedResponse, nil
}

func ConvertWeeklyStatsToMap(weeklyStats map[int]map[string]string) map[string]interface{} {
	converted := make(map[string]interface{})
	for week, stats := range weeklyStats {
		converted[fmt.Sprintf("%d", week)] = stats
	}
	return converted
}

func StructToMap(input interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize struct: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize struct to map: %w", err)
	}

	return result, nil
}
