package schemas

import (
	"github.com/Adejare77/blogPlatform/internals/utilities"
	"gorm.io/gorm"
)

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"not null" binding:"required"`
	Email    string    `gorm:"unique;not null" binding:"required,email"`
	Password string    `gorm:"not null" binding:"required"`
	Posts    []Post    `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Likes    []Like    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (newUser *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUser.Password, err = utilities.HashPassword(newUser.Password)
	return
}
