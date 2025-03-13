package schemas

import "time"

type Like struct {
	AuthorID     uint   `gorm:"not null"`
	UserID       uint   `gorm:"primaryKey;index;uniqueIndex:user_like"`
	LikeableID   string `gorm:"type:text;not null;index:idx_likeable;uniqueIndex:user_like"`
	LikeableType string `gorm:"not null;index:idx_likeable"`
	CreatedAt    time.Time
}
