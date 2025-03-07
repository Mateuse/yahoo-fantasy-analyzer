package repositories

import (
	"fmt"
	"time"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveScheduleGameInDB(scheduleGame *models.ScheduleGame) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Match on the primary key `id`
		DoNothing: true,                          // Ignore the duplicate
	}).Create(scheduleGame).Error

	if err != nil {
		return fmt.Errorf("failed to insert or update data into DB: %w", err)
	}
	return nil
}

func GetTeamNextGameDB(teamAbbrev string) (*models.ScheduleGame, error) {
	var nextGame models.ScheduleGame

	err := DB.Where("(home_team_abbrev = ? OR away_team_abbrev = ?) AND start_time_utc > ?", teamAbbrev, teamAbbrev, time.Now()).
		Order("start_time_utc ASC").
		First(&nextGame).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no upcoming games found for team %s", teamAbbrev)
		}
		return nil, fmt.Errorf("failed to query next game: %s", err)
	}

	return &nextGame, nil
}

func GetTeamNextGameDate(team_abbrev string) (*time.Time, error) {

	nextGame, err := GetTeamNextGameDB(team_abbrev)
	if err != nil {
		return nil, err
	}

	return &nextGame.StartTimeUTC, nil
}
