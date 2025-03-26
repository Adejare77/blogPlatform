package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/gorm"
)

func TotalPosts() (int64, error) {
	var count int64
	if err := config.DB.Model(&schemas.Post{}).
		Where("status = ?", "published").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CreatePost(post *schemas.Post) error {
	return config.DB.Create(&post).Error
}

func FindAllPosts() ([]map[string]any, error) {
	var allPosts []map[string]any

	if err := config.DB.Model(&schemas.Post{}).
		Where("posts.status = ?", "published").
		Select(`
		posts.id AS post_id, 
		users.name AS author_name, 
		posts.title AS post_title, 
		CONCAT(LEFT(posts.content, 150), '...') AS content_excerpt, 
		COUNT(likes.likeable_id) AS likes,
		COUNT(comments.id) AS comments_counts
		`).
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Joins("INNER JOIN users ON users.id = posts.author_id").
		Group("posts.id, users.name").
		Scan(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func FindUserPosts(userID uint, status string) ([]map[string]any, error) {
	var allPosts []map[string]any

	// Check for forbidden access later
	if err := config.DB.Model(&schemas.Post{}).
		Where("posts.author_id = ? AND status = ?", userID, status).
		Select(`
		posts.id AS post_id, 
		users.name AS author_name, 
		posts.title AS post_title, 
		CONCAT(LEFT(posts.content, 150), '...') AS content_excerpt, 
		COUNT(likes.likeable_id) AS likes,
		COUNT(comments.id) AS comments_counts
		`).
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Joins("INNER JOIN users ON users.id = posts.author_id").
		Group("posts.id, users.name").
		Scan(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func FindByPostID(userID uint, postID string, status string) (map[string]any, error) {
	var post map[string]any
	var query *gorm.DB

	if status == "published" {
		query = config.DB.Model(&schemas.Post{}).
			Where("posts.id = ? AND status = ?", postID, status)
	} else { // Authorized Users only
		query = config.DB.Model(&schemas.Post{}).
			Where("posts.id = ? AND author_id = ? AND status = ?", postID, userID, status)
	}

	if err := query.
		Select(`
		posts.id AS post_id, 
		users.name AS author_name,
		title AS post_title,
		content, 
		COUNT(likes.likeable_id) AS likes,
		COUNT(comments.id) AS comments_counts
		`).
		Joins("LEFT JOIN likes ON likes.likeable_id = posts.id").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Joins("INNER JOIN users ON users.id = posts.author_id").
		Group("posts.id, users.name").
		Scan(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func FindPostAuthorID(postID string) (*uint, error) {
	var post schemas.Post

	if err := config.DB.
		Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &post.AuthorID, nil
}

func UpdateUserPost(AuthorID uint, postID string, data map[string]any) error {
	var post schemas.Post

	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not found")
		}
		return err
	}

	if post.AuthorID != AuthorID {
		return errors.New("forbidden")
	}

	if err := config.DB.Model(&schemas.Post{}).
		Where("id = ?", postID).
		Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func DeleteUserPost(userID uint, postID string) error {
	var post schemas.Post
	cursor := config.DB.First(&post, "id = ?", postID)

	if errors.Is(cursor.Error, gorm.ErrRecordNotFound) {
		return errors.New("not found")
	}

	if cursor.Error != nil {
		return cursor.Error
	}

	if post.AuthorID != userID {
		return errors.New("forbidden")
	}

	if err := cursor.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}
