package models

type FantasyPlayerStats struct {
	Stats []StatModifier `gorm:"-"`
}

type ProjectedVsActualStats struct {
	ProjectedStats string
	ActualStats    string
}
