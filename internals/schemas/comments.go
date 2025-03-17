package schemas

import (
	"time"
)

type Comment struct {
	ID        string `gorm:"primaryKey"`
	PostID    string `gorm:"type:uuid;not null"`
	UserID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text" binding:"required"`
	CreatedAt time.Time
	ParentID  *string   `gorm:"type:uuid;index;default:null"`
	Likes     []Like    `gorm:"type:uuid;polymorphic:Likeable;polymorphicValue:Comment;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Replies   []Comment `gorm:"foreignKey:ParentID"`
}
