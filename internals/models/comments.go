package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/gorm"
)

func CreateComment(comment schemas.Comment) error {
	return config.DB.Create(&comment).Error
}

func FindCommentAuthorID(postID string) (*uint, error) {
	var comment schemas.Comment

	if err := config.DB.
		Where("id = ?", postID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &comment.AuthorID, nil
}

func FindCommentsByPostID(postID string) ([]map[string]any, error) {
	var allComments []map[string]any

	err := config.DB.Model(&schemas.Comment{}).
		Where("post_id = ?", postID).
		Select(`
		comments.id AS comment_id, 
		comment_author.name AS comment_author,
		CONCAT(LEFT(comments.content, 100), '...') AS comment_excerpt,
		post_author.name AS post_author,
		COUNT(likes.likeable_id) AS likes
		`).
		Joins("INNER JOIN users AS comment_author ON comment_author.id = comments.author_id").
		Joins("INNER JOIN posts ON posts.id = comments.post_id").
		Joins("INNER JOIN users AS post_author ON post_author.id = posts.author_id").
		Joins("LEFT JOIN likes ON likes.likeable_id = comments.id").
		Group(`comment_id, comment_author, post_author`).
		Scan(&allComments).Error

	if err != nil {
		return nil, err
	}

	return allComments, nil
}

func UpdateComment(userID uint, filter schemas.CommentUriParam, data schemas.CommentBody) error {
	cursor := config.DB.
		Where("author_id = ? AND id = ?", userID, filter.CommentID).
		Update("content", data.Content)

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func DeleteComment(userID uint, comment schemas.CommentUriParam) error {
	cursor := config.DB.
		Where("author_id = ? AND id = ?", userID, comment.CommentID).
		Delete(&schemas.Comment{})

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}
