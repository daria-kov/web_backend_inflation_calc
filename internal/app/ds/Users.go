package ds

import "time"

type Users struct {
	UserID    uint      `gorm:"primaryKey"`
	FullName  string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time `gorm:"not null"`
}
