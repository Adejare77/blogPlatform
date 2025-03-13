package models

import (
	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func CreateUser(user *schemas.User) error {
	return config.DB.Create(&user).Error
}

func GetUserInfo(email string) (schemas.User, error) {
	user := schemas.User{}
	if err := config.DB.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return schemas.User{}, err
	}

	return user, nil
}

// func GetUserByPostID(postID string) (schemas.User, error) {
// 	var post schemas.Post
// 	cursor := config.DB.Preload("User").First(&post, "id = ?", postID)
// 	if err := cursor.Error; err != nil {
// 		return schemas.User{}, err
// 	}
// 	return post.User, nil
// }
