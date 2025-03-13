package models

import (
	"errors"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/gorm"
)

func CreatePost(post *schemas.Post) error {
	return config.DB.Create(&post).Error
}

func GetAllPosts() ([]map[string]any, error) {
	var allPosts []map[string]any

	if err := config.DB.Model(&schemas.Post{}).
		Select("posts.id, users.name, posts.title, CONCAT(LEFT(posts.content, 150), ...) AS content, COUNT(likes.id) AS likes").
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Joins("LEFT JOIN users ON users.id = posts.author_id").
		Group("posts.id").
		Scan(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func GetPosts(userID uint) ([]map[string]any, error) {
	var allPosts []map[string]any
	if err := config.DB.Model(&schemas.Post{}).
		Where("posts.author_id = ?", userID).
		Select("posts.id, users.name, posts.title, CONCAT(LEFT(posts.content, 150), ...) AS content, COUNT(likes.id) AS likes").
		Joins("LEFT JOIN likes ON posts.id = likes.likeable_id AND likes.likeable_type = ?", "Post").
		Joins("LEFT JOIN users ON users.id = posts.author_id").
		Group("posts.id").
		Scan(&allPosts).Error; err != nil {
		return nil, err
	}

	return allPosts, nil
}

func GetPostByID(postID string) (map[string]any, error) {
	var post map[string]any

	if err := config.DB.Model(&schemas.Post{}).
		Where("posts.id = ?", postID).
		Select("posts.id, posts.author_id, title, content, COUNT(*) AS likes").
		Joins("LEFT JOIN likes ON likes.likeable_id = posts.id").
		Scan(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func GetAuthorIDByPostID(postID string) (uint, error) {
	var authorID uint
	if err := config.DB.
		Where("id = ?", postID).
		Select("author_id").Scan(&authorID).Error; err != nil {
		return 0, err
	}
	return authorID, nil
}

func UpdatePost(data map[string]any) error {
	delete(data, "id") // deletes possibility of updating "postID"
	if err := config.DB.
		Model(&schemas.Post{}).
		Where("author_id = ? AND id = ?", data["userID"], data["postID"]).
		Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func DeletePost(userID uint, postID string) error {
	var post schemas.Post
	cursor := config.DB.First(&post, "id = ?", postID)

	if errors.Is(cursor.Error, gorm.ErrRecordNotFound) {
		return errors.New("record Not Found")
	} else if post.AuthorID != userID {
		return errors.New("only Authorized Personnel is Allowed")
	} else {
		if err := config.DB.Delete(&post).Error; err != nil {
			return err
		}
	}
	return nil
}
