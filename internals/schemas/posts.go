package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        string `gorm:"primaryKey"`
	AuthorID  uint   `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title     string `binding:"required"`
	Content   string `gorm:"type:text" binding:"required"`
	Status    string `gorm:"default:'draft'" binding:"required,oneof=unpublished published"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     []Like    `gorm:"polymorphic:Likeable;polymorphicValue:Post;constraint:OnDelete:CASCADE"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE" binding:"-"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	post.ID = uuid.New().String()
	return
}
