package models

import "time"

type Users struct {
	Count int    `xml:"count,attr"`
	User  []User `xml:"user"`
}

type User struct {
	GUID  string `xml:"guid"`
	Games Games  `xml:"games"`
}

type Games struct {
	Count int    `xml:"count,attr"`
	Game  []Game `xml:"game"`
}

type Game struct {
	GameKey            string  `xml:"game_key"`
	GameID             string  `xml:"game_id"`
	Name               string  `xml:"name"`
	Code               string  `xml:"code"`
	Type               string  `xml:"type"`
	URL                string  `xml:"url"`
	Season             string  `xml:"season"`
	IsRegistrationOver int     `xml:"is_registration_over"`
	IsGameOver         int     `xml:"is_game_over"`
	IsOffSeason        int     `xml:"is_offseason"`
	Leagues            Leagues `xml:"leagues"`
}

type Leagues struct {
	Count  int      `xml:"count,attr"`
	League []League `xml:"league"`
}

type FantasyContent struct {
	League League `xml:"league"`
}

type RosterPosition struct {
	Position           string `gorm:"column:position"`
	PositionType       string `gorm:"column:position_type"`
	Count              int    `gorm:"column:count"`
	IsStartingPosition bool   `gorm:"column:is_starting_position"`
}

type StatModifier struct {
	StatID   string  `gorm:"column:stat_id"`
	Value    float64 `gorm:"column:value"`
	StatName string  `gorm:"column:name`
}

type League struct {
	LeagueID              string           `gorm:"primaryKey"`
	LeagueKey             string           `gorm:"column:league_key"`
	Name                  string           `gorm:"column:name"`
	URL                   string           `gorm:"column:url"`
	LogoURL               string           `gorm:"column:logo_url"`
	DraftStatus           string           `gorm:"column:draft_status"`
	NumTeams              int              `gorm:"column:num_teams"`
	EditKey               string           `gorm:"column:edit_key"`
	WeeklyDeadline        string           `gorm:"column:weekly_deadline"`
	LeagueUpdateTimestamp int64            `gorm:"column:league_update_timestamp"`
	ScoringType           string           `gorm:"column:scoring_type"`
	LeagueType            string           `gorm:"column:league_type"`
	FeloTier              string           `gorm:"column:felo_tier"`
	AllowAddToDLExtraPos  bool             `gorm:"column:allow_add_to_dl_extra_pos"`
	IsProLeague           bool             `gorm:"column:is_pro_league"`
	IsCashLeague          bool             `gorm:"column:is_cash_league"`
	CurrentWeek           int              `gorm:"column:current_week"`
	StartWeek             int              `gorm:"column:start_week"`
	StartDate             time.Time        `gorm:"column:start_date"`
	EndWeek               int              `gorm:"column:end_week"`
	EndDate               time.Time        `gorm:"column:end_date"`
	IsPlusLeague          bool             `gorm:"column:is_plus_league"`
	GameCode              string           `gorm:"column:game_code"`
	Season                string           `gorm:"column:season"`
	MaxTeams              int              `gorm:"column:max_teams"`
	RosterPositions       []RosterPosition `gorm:"-"`
	StatModifiers         []StatModifier   `gorm:"-"`
	LastUpdated           time.Time        `gorm:"autoUpdateTime"`
}

type Standings struct {
	Teams []Team `xml:"teams>team"`
}

type Team struct {
	TeamKey           string `xml:"team_key"`
	TeamID            string `xml:"team_id"`
	Name              string `xml:"name"`
	URL               string `xml:"url"`
	LogoURL           string `xml:"team_logos>team_logo>url"`
	WaiverPriority    int    `xml:"waiver_priority"`
	NumberOfMoves     int    `xml:"number_of_moves"`
	NumberOfTrades    int    `xml:"number_of_trades"`
	LeagueScoringType string `xml:"league_scoring_type"`
	DraftPosition     int    `xml:"draft_position"`
	CurrentWeekStats  string `xml:"team_stats>week"`
	CurrentWeekPoints string `xml:"team_points>total"`
	ProjectedPoints   string `xml:"team_projected_points>total"`
	FinalPoints       string `xml:"team_live_projected_points>total"`
	RemainingGames    int    `xml:"team_remaining_games>total>remaining_games"`
	CompletedGames    int    `xml:"team_remaining_games>total>completed_games"`
}

type Player struct {
	PlayerID           string         `gorm:"column:player_id"`            // Unique player ID
	PlayerKey          string         `gorm:"column:player_key"`           // Unique player key
	Name               PlayerName     `gorm:"embedded"`                    // Embedded struct for player name details
	TeamFullName       string         `gorm:"column:team_full_name"`       // Full team name
	TeamAbbreviation   string         `gorm:"column:team_abbr"`            // Team abbreviation
	TeamURL            string         `gorm:"column:team_url"`             // Team's URL
	UniformNumber      string         `gorm:"column:uniform_number"`       // Player's jersey number
	DisplayPosition    string         `gorm:"column:display_position"`     // Displayed position (e.g., C, LW, RW, D)
	HeadshotURL        string         `gorm:"column:headshot_url"`         // URL to the player's headshot
	ImageURL           string         `gorm:"column:image_url"`            // URL to the player's full image
	IsUndroppable      bool           `gorm:"column:is_undroppable"`       // Indicates if the player is undroppable
	PositionType       string         `gorm:"column:position_type"`        // Player position type (e.g., "P" for player)
	EligiblePositions  []string       `gorm:"-"`                           // List of eligible positions
	PlayerNotes        bool           `gorm:"column:player_notes"`         // Indicates if the player has notes
	RecentNotes        bool           `gorm:"column:recent_notes"`         // Indicates if there are recent player notes
	PlayerNotesUpdated int64          `gorm:"column:player_notes_updated"` // Timestamp of the last notes update
	Stats              []Stat         `gorm:"-"`                           // Regular stats
	AdvancedStats      []AdvancedStat `gorm:"-"`                           // Advanced stats
}

type PlayerName struct {
	Full       string `gorm:"column:full_name"`   // Full name
	First      string `gorm:"column:first_name"`  // First name
	Last       string `gorm:"column:last_name"`   // Last name
	AsciiFirst string `gorm:"column:ascii_first"` // ASCII version of the first name
	AsciiLast  string `gorm:"column:ascii_last"`  // ASCII version of the last name
}

type Stat struct {
	StatID string `gorm:"column:stat_id"` // Stat ID
	Value  string `gorm:"column:value"`   // Value of the stat
}

type AdvancedStat struct {
	StatID string `gorm:"column:stat_id"` // Stat ID
	Value  string `gorm:"column:value"`   // Value of the stat
}
