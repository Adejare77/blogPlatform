package schemas

import "time"

type Like struct {
	UserID       uint   `gorm:"index;uniqueIndex:user_like"`
	LikeableID   string `gorm:"type:uuid;not null;index:idx_likeable;uniqueIndex:user_like"`
	LikeableType string `gorm:"not null;index:idx_likeable"`
	CreatedAt    time.Time
}

type LikedQueryParams struct {
	Type string `form:"type" binding:"required,oneof=post comment"`
}
