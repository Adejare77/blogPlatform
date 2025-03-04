package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func PostComment(comment schemas.Comment) error {
	if err := config.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func GetCommentsByPostID(postID string) ([]schemas.Comment, error) {
	var allComments []schemas.Comment

	if err := config.DB.
		Find(&allComments, "post_id = ?", postID).
		Error; err != nil {
		return nil, err
	}

	return allComments, nil
}

func UpdateComment(userID uint, commentID string, updateData map[string]string) error {
	cursor := config.DB.
		Where("user_id = ? AND id = ?", userID, commentID).
		Update("content", updateData["content"])

	if cursor.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if cursor.Error != nil {
		return cursor.Error
	}

	return nil
}

func DeleteComment(userID uint, commentID string) error {
	cursor := config.DB.
		Where("user_id = ? AND id = ?", userID, commentID).
		Delete(&schemas.Comment{})

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

// func ReplyComment(userID uint, postID string, commentID string) {
// 	config.DB.W
// }
