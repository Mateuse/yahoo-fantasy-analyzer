package models

type PlayerDetails struct {
	Player            NHLPlayer
	YahooPlayerID     string
	EligiblePositions []string
	Stats             []Stat
	AdvancedStats     []Stat
}
