package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func GetLikedPosts(userID uint) ([]schemas.Post, error) {
	var likeableIDs []uint

	if err := config.DB.Model(&schemas.Like{}).
		Where("user_id = ? AND likeable_type = ?", "Post").
		Pluck("likeable_id", &likeableIDs).Error; err != nil {
		return nil, err
	}

	var user schemas.User

	if err := config.DB.
		Preload("Posts", "id IN ?", likeableIDs).
		First(&user).Error; err != nil {
		return nil, err
	}

	return user.Posts, nil
}

func LikePostOrComment(userID uint, postID string, parent string) error {
	cursor := config.DB.
		Where("user_id = ? AND likeable_id = ?", userID, postID).
		FirstOrCreate(&schemas.Like{
			UserID:       userID,
			LikeableID:   postID,
			LikeableType: parent,
		})

	return cursor.Error
}

func UnlikePostOrComment(userID uint, postID string, parent string) error {
	result := config.DB.
		Where("user_id = ? AND likeable_id = ? AND likeable_type = ?", userID, postID, parent).
		Delete(&schemas.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}
