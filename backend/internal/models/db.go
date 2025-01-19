package models

type RefreshToken struct {
	UserId       string `gorm:"primaryKey"`
	RefreshToken string `gorm:"size255;not null"`
}
