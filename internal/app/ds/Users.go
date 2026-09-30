package ds

type Users struct {
	UserID   uint   `gorm:"primaryKey"`
	FullName string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
}
