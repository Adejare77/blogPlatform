package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        string `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Title     string `binding:"required"`
	Content   string `gorm:"type:text" binding:"required"`
	Status    string `gorm:"default:'draft'" binding:"omitempty,oneof=draft published"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     []Like    `gorm:"polymorphic:Likeable;polymorphicValue:Post;constraint:OnDelete:CASCADE"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constrait:OnDelete:CASCADE" binding:"-"` // Optional: for eager loading of the associated User (belongs-to relationship)
}

func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	post.ID = uuid.New().String()
	return
}
