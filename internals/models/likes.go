package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func GetLikedPosts(userID uint) ([]map[string]string, error) {
	var result []map[string]string

	if err := config.DB.Model(&schemas.Like{}).
		Select("users.name, posts.title, LEFT(posts.content, 200) AS content").
		Joins("LEFT JOIN posts ON posts.id = likes.likeable_id").
		Joins("LEFT JOIN users ON likes.author_id = users.id").
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func LikePostOrComment(data schemas.Like) error {
	return config.DB.FirstOrCreate(&data).Error
}

func UnlikePostOrComment(like schemas.Like) error {
	cursor := config.DB.Delete(&like)

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}
