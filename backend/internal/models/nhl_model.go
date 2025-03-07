package models

import "time"

type ScheduleGame struct {
	ID             int64     `json:"id"`
	Season         int       `json:"season"`
	GameType       int       `json:"gameType"`
	GameDate       string    `json:"gameDate"`
	StartTimeUTC   time.Time `json:"startTimeUTC"` // Use time.Time for compatibility
	HomeTeamAbbrev string    `json:"homeTeam"`
	AwayTeamAbbrev string    `json:"awayTeam"`
}

type PlayerIDMapping struct {
	YahooPlayerID string
	NHLPlayerID   string
	PlayerName    string
	TeamAbbr      string
}

type NHLRoster struct {
	Forwards   []NHLPlayer `json:"forwards"`
	Defensemen []NHLPlayer `json:"defensemen"`
	Goalies    []NHLPlayer `json:"goalies"`
}

type NHLPlayer struct {
	ID             int    `gorm:"primaryKey;column:id"`
	Headshot       string `gorm:"column:headshot"`
	FirstName      string `gorm:"column:first_name"`
	LastName       string `gorm:"column:last_name"`
	SweaterNumber  int    `gorm:"column:sweater_number"`
	PositionCode   string `gorm:"column:position_code"`
	ShootsCatches  string `gorm:"column:shoots_catches"`
	HeightInInches int    `gorm:"column:height_in_inches"`
	WeightInPounds int    `gorm:"column:weight_in_pounds"`
	HeightInCM     int    `gorm:"column:height_in_cm"`
	WeightInKG     int    `gorm:"column:weight_in_kg"`
	BirthDate      string `gorm:"column:birth_date"`
	BirthCity      string `gorm:"column:birth_city"`
	BirthCountry   string `gorm:"column:birth_country"`
	BirthState     string `gorm:"column:birth_state"`
	Team           string `gorm:"column:team"`
}

type PlayerGameStat struct {
	GameID            string `json:"gameId"`
	PlayerID          string `json:"playerId"`
	TeamAbbrev        string `json:"teamAbbrev"`
	HomeRoadFlag      string `json:"homeRoadFlag"`
	GameDate          string `json:"gameDate"`
	Goals             int    `json:"goals"`
	Assists           int    `json:"assists"`
	Team              string `json:"team"`
	Opponent          string `json:"opponent"`
	Points            int    `json:"points"`
	PlusMinus         int    `json:"plusMinus"`
	PowerPlayGoals    int    `json:"powerPlayGoals"`
	PowerPlayPoints   int    `json:"powerPlayPoints"`
	GameWinningGoals  int    `json:"gameWinningGoals"`
	OTGoals           int    `json:"otGoals"`
	Shots             int    `json:"shots"`
	Shifts            int    `json:"shifts"`
	ShorthandedGoals  int    `json:"shorthandedGoals"`
	ShorthandedPoints int    `json:"shorthandedPoints"`
	OpponentAbbrev    string `json:"opponentAbbrev"`
	PIM               int    `json:"pim"`
	TOI               string `json:"toi"`
}
