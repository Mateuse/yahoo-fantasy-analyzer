package repositories

import (
	"fmt"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveFTeamMatchups(matchups []*models.Matchup) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&matchups).Error; err != nil {
			return fmt.Errorf("failed to insert matchups: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
