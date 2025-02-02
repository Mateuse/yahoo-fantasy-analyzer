package models

type FantasyPlayerStats struct {
	Stats []StatModifier `gorm:"-"`
}
