package ds

import "time"

type Like struct {
	LikeID     uint      `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"not null"`
	UserID     uint      `gorm:"not null"`
	CategoryID uint      `gorm:"not null"`
}
