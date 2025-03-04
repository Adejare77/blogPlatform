package models

import (
	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func CreateUser(user *schemas.User) error {
	result := config.DB.Create(user)
	return result.Error
}

func GetUserInfo(email string) (*schemas.User, bool) {
	user := &schemas.User{}
	cursor := config.DB.
		Model(&schemas.User{}).
		Where("email = ?", email).
		First(user)

	if cursor.RowsAffected == 0 {
		return nil, false
	}
	return user, true
}

func GetUserByPostID(postID string) (schemas.User, error) {
	var post schemas.Post
	cursor := config.DB.Preload("User").First(&post, "id = ?", postID)
	if err := cursor.Error; err != nil {
		return schemas.User{}, err
	}
	return post.User, nil
}
