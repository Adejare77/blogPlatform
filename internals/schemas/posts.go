package schemas

import (
	"time"
)

type Post struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID  uint      `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title     string    `binding:"required"`
	Content   string    `gorm:"type:text" binding:"required"`
	Status    string    `gorm:"default:'draft';index" binding:"omitempty,oneof=draft published"`
	Likes     []Like    `gorm:"polymorphic:Likeable;polymorphicValue:Post;constraint:OnDelete:CASCADE"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

type UserPostsStats struct {
	TotalPosts  int64 `json:"total_user_posts"`
	TotalDrafts int64 `json:"total_user_drafts"`
}

type PostURIParams struct {
	PostID string `uri:"post_id" binding:"required,uuid"`
}

type StatusQueryParams struct {
	Status string `form:"status" binding:"required,oneof=draft published"`
}

type FilterQueryParams struct {
	Page  int `form:"page" binding:"numeric,min=1"`
	Limit int `form:"limit" binding:"numeric,min=1"`
}
