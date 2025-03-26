package models

import (
	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
)

func FindLikesByUser(userID uint, targetType string) ([]map[string]string, error) {
	var result []map[string]string

	query := config.DB.Model(&schemas.Like{}).
		Where("user_id = ? AND likeable_type = ?", userID, targetType)

	if targetType == "post" {
		err := query.Select(`
		posts.id AS post_id,
		posts.title AS post_title,
		CONCAT(LEFT(posts.content, 200), '...') AS post_excerpt,
		posts.created_at AS post_created_at,
		users.name AS post_author,
		`).
			Joins(`INNER JOIN posts ON posts.id = likes.likeable_id`).
			Joins(`INNER JOIN users ON likes.author_id = users.id`).
			Find(&result).Error

		if err != nil {
			return nil, err
		}

	} else {
		err := query.Select(`
		comments.id AS comment_id,
		CONCAT(LEFT(comments.content, 100), '...') AS comment_excerpt,
		comment_user.name AS comment_author,
		comments.created_at AS comment_created_at
		posts.title AS parent_post_title,
		posts.id AS parent_post_id,
		post_user.name AS post_author,
		`).
			Joins(`INNER JOIN comments ON comments_id = likes.likeable_ids`).
			Joins(`INNER JOIN posts ON posts.id = comments.post_id`).
			Joins(`INNER JOIN users AS comment_user ON comments.author_id = users.id`).
			Joins(`INNER JOIN users AS post_user ON posts.author_id = users.id`).
			Find(&result).Error

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func CreateLike(data schemas.Like) error {
	return config.DB.Create(&data).Error
}

func DeleteLike(userID uint, targetID string) error {
	return config.DB.
		Where("user_id = ? AND likeableID = ?", userID, targetID).
		Delete(&schemas.Like{}).Error
}
