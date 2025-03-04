package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        string `gorm:"primaryKey"`
	PostID    string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text" binding:"required"`
	CreatedAt time.Time
	ParentID  *string   `gorm:"index;default:null"`
	Likes     []Like    `gorm:"polymorphic:Likeable;polymorphicValue:Comment;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Replies   []Comment `gorm:"foreignKey:ParentID"`
}

func (comment *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	comment.ID = uuid.New().String()
	return
}
