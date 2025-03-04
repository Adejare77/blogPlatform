package models

import (
	"errors"
	"net/http"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/gorm"
)

func CreatePost(post *schemas.Post) error {
	result := config.DB.Create(post)

	return result.Error
}

func GetAllPosts() ([]map[string]interface{}, error) {
	var allPosts []map[string]interface{}

	if err := config.DB.Model(&schemas.Post{}).
		Select("posts.id, posts.title, posts.content, COUNT(likes.id) AS count").
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Group("posts.id").
		Find(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func GetPosts(userID uint) ([]map[string]interface{}, error) {
	var allPosts []map[string]interface{}
	if err := config.DB.Model(&schemas.Post{}).
		Select("posts.id, posts.title, posts.content, COUNT(likes.id) AS count").
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Group("posts.id").
		Find(&allPosts, "posts.user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func GetPost(postID string) (map[string]interface{}, error) {
	var post map[string]interface{}
	var count int64

	if err := config.DB.
		Preload("Likes", "likeable_id = ?", postID, func(db *gorm.DB) *gorm.DB {
			return db.Count(&count)
		}).
		Where("id = ?", postID).
		Select("id, title, content").
		First(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func UpdatePost(userID uint, postID string, data map[string]interface{}) error {
	delete(data, "id") // deletes possibility of updating "postID"
	if err := config.DB.
		Model(&schemas.Post{}).
		Where("user_id = ? AND id = ?", userID, postID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DeletePost(userID uint, postID string) (int, error) {
	var post schemas.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	if post.UserID != userID {
		return http.StatusForbidden, errors.New("only Authorized Personnel is Allowed")
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
