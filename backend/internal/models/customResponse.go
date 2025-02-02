package models

type CustomResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type PlayerRanksResponse struct {
	PlayerID    string       `json:"player_id"`    // Player ID
	PlayerRanks []PlayerRank `json:"player_ranks"` // Player ranks
}
