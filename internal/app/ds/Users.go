package ds

type Users struct {
	UserID   uint   `gorm:"primaryKey"`
	FullName string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}
