package schemas

import (
	"time"
)

type Comment struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PostID    string    `gorm:"type:uuid;not null"`
	AuthorID  uint      `gorm:"not null"`
	Content   string    `gorm:"type:text"`
	ParentID  *string   `gorm:"type:uuid;index;default:null"`
	Likes     []Like    `gorm:"type:uuid;polymorphic:Likeable;polymorphicValue:Comment;constraint:OnDelete:CASCADE"`
	Replies   []Comment `gorm:"foreignKey:ParentID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentUriParam struct {
	PostID    string  `uri:"post_id" binding:"required,uuid"`
	CommentID string  `uri:"comment_id" binding:"required,uuid"`
	ParentID  *string `uri:"parent_id" binding:"omitempty,uuid"`
}

type CommentBody struct {
	Content string `json:"content" binding:"required"`
}
