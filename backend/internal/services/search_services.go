package services

import "github.com/mateuse/yahoo-fantasy-analyzer/internal/models"

func GetPlayerByName(userSession, playerName string) models.Player
