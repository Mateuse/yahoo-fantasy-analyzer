package repositories

import (
	"fmt"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveStatWinnerWeeklyMatchup(statWinnerWeeklyMatchup []*models.StatWinnerWeeklyMatchup) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&statWinnerWeeklyMatchup).Error; err != nil {
			return fmt.Errorf("failed to insert stat winner weekly matchup: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
