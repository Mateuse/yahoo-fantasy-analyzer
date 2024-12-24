package responses

import (
	"encoding/xml"
	"fmt"
)

// Define structs to represent the XML response
type FantasyContent struct {
	XMLName xml.Name `xml:"fantasy_content"`
	Users   Users    `xml:"users"`
}

type Users struct {
	User User `xml:"user"`
}

type User struct {
	GUID string `xml:"guid"`
}

// ParseFantasyContent parses the XML response and extracts the GUID
func ParseFantasyContent(body []byte) (string, error) {
	var fantasyContent FantasyContent

	// Unmarshal the XML into the FantasyContent struct
	err := xml.Unmarshal(body, &fantasyContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse XML: %w", err)
	}

	// Check if the GUID exists
	if fantasyContent.Users.User.GUID == "" {
		return "", fmt.Errorf("user ID (GUID) not found in response")
	}

	return fantasyContent.Users.User.GUID, nil
}
