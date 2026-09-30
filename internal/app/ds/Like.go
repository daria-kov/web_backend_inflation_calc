package ds

type Like struct {
	LikeID     uint `gorm:"primaryKey"`
	UserID     uint `gorm:"not null"`
	CategoryID uint `gorm:"not null"`
}
