package repositories

import (
	"fmt"
	"log"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
)

func SaveLeagueTeamsToDB(leagueTeams []models.LeagueTeam) error {

	err := DB.Create(&leagueTeams).Error
	if err != nil {
		return fmt.Errorf("failed to insert league teams: %w", err)
	}

	log.Printf("Sucessfully inserted %d league teams", len(leagueTeams))
	return nil
}

// GetAllLeagueTeamsFromDB fetches all league teams from the database
func GetAllLeagueTeamsFromDB(leagueId string) ([]models.LeagueTeam, error) {
	var leagueTeams []models.LeagueTeam

	// Query teams where league_id matches the input
	err := DB.Where("league_id = ?", leagueId).Find(&leagueTeams).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teams for league %s: %w", leagueId, err)
	}

	if len(leagueTeams) == 0 {
		return []models.LeagueTeam{}, nil
	}

	return leagueTeams, nil
}

func SaveTeamMatchups(matchups []*models.Matchup, teamWeeklyStats []*models.TeamWeeklyStats, statWinnerWeeklyMatchup []*models.StatWinnerWeeklyMatchup) error {
	if err := SaveFTeamMatchups(matchups); err != nil {
		return fmt.Errorf("failed to save matchups: %w", err)
	}

	if err := SaveTeamWeeklyStats(teamWeeklyStats); err != nil {
		return fmt.Errorf("failed to save team weekly stats: %w", err)
	}

	if err := SaveStatWinnerWeeklyMatchup(statWinnerWeeklyMatchup); err != nil {
		return fmt.Errorf("failed to save stat winner weekly matchups: %w", err)
	}

	return nil
}
