package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/gorm"
)

func PublishedPosts(db *gorm.DB) *gorm.DB {
	return db.Where("posts.status = ?", "published")
}

func CreatePost(post *schemas.Post) error {
	return config.DB.Create(&post).Error
}

func FindAllPosts(page int, limit int) ([]map[string]any, *int64, error) {
	var allPosts []map[string]any
	offset := (page - 1) * limit

	query := config.DB.Model(&schemas.Post{}).
		Scopes(PublishedPosts).
		Select(`
		posts.id AS post_id, 
		users.name AS author_name, 
		posts.title AS post_title, 
		CONCAT(LEFT(posts.content, 150), '...') AS content_excerpt, 
		COUNT(DISTINCT likes.likeable_id) AS likes,
		COUNT(DISTINCT comments.id) AS comments_counts
		`).
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "post").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Joins("INNER JOIN users ON users.id = posts.author_id").
		Group("posts.id, users.name").
		Order("posts.created_at DESC").
		Offset(offset).
		Limit(limit)

	if query.Scan(&allPosts).Error != nil {
		return nil, nil, query.Error
	}

	var totalPosts int64

	if err := config.DB.Model(&schemas.Post{}).
		Scopes(PublishedPosts).
		Count(&totalPosts).Error; err != nil {
		return nil, nil, err
	}

	return allPosts, &totalPosts, nil
}

func FindUserPosts(userID uint, status string, page int, limit int) ([]map[string]any, *schemas.UserPostsStats, error) {
	var allPosts []map[string]any
	offset := (page - 1) * limit

	if err := config.DB.Model(&schemas.Post{}).
		Where("posts.author_id = ? AND status = ?", userID, status).
		Select(`
		posts.id AS post_id, 
		users.name AS author_name, 
		posts.title AS post_title, 
		CONCAT(LEFT(posts.content, 150), '...') AS content_excerpt, 
		COUNT(DISTINCT likes.likeable_id) AS likes,
		COUNT(DISTINCT comments.id) AS comments_counts
		`).
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "post").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Joins("INNER JOIN users ON users.id = posts.author_id").
		Group("posts.id, users.name").
		Order("posts.created_at DESC").
		Offset(offset).Limit(limit).
		Scan(&allPosts).Error; err != nil {
		return nil, nil, err
	}

	var userStats schemas.UserPostsStats

	if err := config.DB.Model(&schemas.Post{}).
		Where("author_id = ?", userID).
		Select(`
		SUM(CASE WHEN status = 'published' THEN 1 ELSE 0 END) AS total_posts,
		SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END) AS total_drafts
		`).Scan(&userStats).Error; err != nil {
		return nil, nil, err
	}

	return allPosts, &userStats, nil
}

func FindByPostID(userID uint, postID string, status string) (map[string]any, error) {
	var post map[string]any
	var query *gorm.DB

	if status == "published" {
		query = config.DB.Model(&schemas.Post{}).
			Where("posts.id = ? AND status = ?", postID, status)
	} else { // Authorized Users only
		query = config.DB.Model(&schemas.Post{}).
			Where("posts.id = ? AND posts.author_id = ? AND status = ?", postID, userID, status)
	}

	if err := query.
		Select(`
		posts.id AS post_id, 
		users.name AS author_name,
		title AS post_title,
		posts.content AS content, 
		COUNT(DISTINCT likes.likeable_id) AS likes,
		COUNT(DISTINCT comments.id) AS comments_counts
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
	cursor := config.DB.Model(&schemas.Post{}).
		Where("id = ? AND author_id = ?", postID, AuthorID).
		Updates(data)

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found or unauthorized")
	}

	return nil
}

func DeleteUserPost(AuthorID uint, postID string) error {
	cursor := config.DB.Where("id = ?  AND author_id = ?", postID, AuthorID).
		Delete(&schemas.Post{})

	if cursor.Error != nil {
		return cursor.Error
	}

	if cursor.RowsAffected == 0 {
		return errors.New("not found or unauthorized")
	}

	return nil
}
