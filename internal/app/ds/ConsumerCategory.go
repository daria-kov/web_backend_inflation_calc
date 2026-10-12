package ds

import "time"

type ConsumerCategory struct {
	CategoryID     uint       `gorm:"primaryKey"`
	Title          string     `gorm:"type:varchar(100);not null"`
	Description    string     `gorm:"type:text"`
	Status         string     `gorm:"type:varchar(20);not null"`
	ImageURL       string     `gorm:"type:varchar(255)"`
	VideoURL       string     `gorm:"type:varchar(255)"`
	BasketWeight   float64    `gorm:"not null"`
	HistoricalDate time.Time  `gorm:"not null"`
	CreatedAt      time.Time  `gorm:"not null"`
	PublishedAt    *time.Time `gorm:""`
	CreatorID      uint       `gorm:"not null"`

	LikesCount int `gorm:"->"`
	IsLiked    int `gorm:"->;-:migration"`
}
