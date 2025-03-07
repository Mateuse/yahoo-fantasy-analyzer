package repositories

import (
	"fmt"
	"log"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
)

func AddRefreshToken(userId, refreshToken string) error {
	refreshTokenEntry := models.RefreshToken{UserId: userId, RefreshToken: refreshToken}
	if err := DB.Create(&refreshTokenEntry).Error; err != nil {
		return fmt.Errorf("failed to add refresh token for user: %w", err)
	}
	log.Printf("User %s Refresh Token added", userId)
	return nil
}

func GetRefreshToken(userId string) (string, error) {
	var refreshTokenEntry models.RefreshToken

	if err := DB.First(&refreshTokenEntry, "user_id = ?", userId).Error; err != nil {
		return "", fmt.Errorf("failed to get refresh token for user %s: $w", userId, err)
	}

	return refreshTokenEntry.RefreshToken, nil
}
