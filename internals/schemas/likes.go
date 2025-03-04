package schemas

import "time"

type Like struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	LikeableID   string `gorm:"type:text;not null;index:idx_likeable"`
	LikeableType string `gorm:"not null;index:idx_likeable"`
	CreatedAt    time.Time
	User         User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Belongs to User
}
