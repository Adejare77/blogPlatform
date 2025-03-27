package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreatePost(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)

	var post schemas.Post

	if err := ctx.ShouldBind(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	post.AuthorID = userID

	if err := models.CreatePost(&post); err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	if err := config.IncrementTotalPosts(); err != nil {
		logrus.WithFields(logrus.Fields{
			"Warning": "Failed to Increase Total blog Posts",
		})
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "post created successfully",
	})
}

func GetAllPosts(ctx *gin.Context) {
	var filters schemas.FilterQueryParams

	filters.Page = 1
	filters.Limit = 20

	if err := ctx.ShouldBindQuery(&filters); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	posts, err := models.FindAllPosts(filters.Page, filters.Limit)
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	} else if len(posts) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no record found",
		})
		return
	}

	prev := fmt.Sprintf("/index?page=%d&limit=%d", filters.Page-1, filters.Limit)
	next := fmt.Sprintf("/index?page=%d&limit=%d", filters.Page+1, filters.Limit)

	if filters.Page == 1 {
		prev = "null"
	}
	if filters.Page == int(config.TotalPosts) {
		next = "null"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"page":        filters.Page,
			"limit":       filters.Limit,
			"total_posts": config.TotalPosts,
			"links": gin.H{
				"next": next,
				"prev": prev,
			},
		},
	})
}

func GetUserPosts(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(uint)
	var status schemas.StatusQueryParams
	var filters schemas.FilterQueryParams

	filters.Limit = 20 // Default limit
	filters.Page = 1   // Default page
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	status.Status = "published" // Default if not given
	if err := ctx.ShouldBindQuery(&status); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	posts, err := models.FindUserPosts(currentUser, status.Status, filters.Page, filters.Limit)
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	if len(posts) == 0 {
		ctx.JSON(http.StatusOK, []map[string]any{})
		return
	}

	// ctx.JSON(http.StatusOK, posts)
	ctx.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"page":             filters.Page,
			"limit":            filters.Page,
			"total_user_posts": "",
			"links":            "",
		},
	})
}

func GetPostByID(ctx *gin.Context) {
	var post schemas.PostURIParams
	var status schemas.StatusQueryParams
	var currentUser uint // Defaults to 0 (No user)

	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	status.Status = "published" // Defaults to published

	if err := ctx.ShouldBindQuery(&status); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	if status.Status == "draft" {
		userID, err := getCurrentUser(ctx)
		if err != nil {
			handlers.Unauthorized(ctx, "login required", "Access to Unauthorized files")
			return
		}
		currentUser = userID
	}

	result, err := models.FindByPostID(currentUser, post.PostID, status.Status)
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "no record found",
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func UpdatePost(ctx *gin.Context) {
	AuthorID := ctx.MustGet("currentUser").(uint)
	var post schemas.PostURIParams

	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	type PostUpdateDTO struct {
		Title   *string `json:"title" binding:"omitempty"`
		Content *string `json:"content" binding:"omitempty"`
		Status  *string `json:"status" binding:"omitempty,oneof=draft published"`
	}

	var dto PostUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	updateData := make(map[string]any)
	if dto.Content != nil {
		updateData["content"] = dto.Content
	}

	if dto.Title != nil {
		updateData["title"] = dto.Title
	}

	if dto.Status != nil {
		updateData["status"] = dto.Status
	}

	if err := models.UpdateUserPost(AuthorID, post.PostID, updateData); err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "Record not Found", err.Error())
			return
		}
		if strings.Contains(err.Error(), "forbidden") {
			handlers.Forbidden(ctx, "Forbidden", err.Error())
			return
		}

		handlers.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "updated successfully",
	})
}

func DeletePost(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	var post schemas.PostURIParams

	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	if err := models.DeleteUserPost(userID, post.PostID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "Record not Found", err.Error())
		} else if strings.Contains(err.Error(), "forbidden") {
			handlers.Forbidden(ctx, "Forbidden", err.Error())
		} else {
			handlers.InternalServerError(ctx, err.Error())
		}
		return
	}

	if err := config.DecrementTotalPosts(); err != nil {
		logrus.WithField("Warning", "Failed to Decrease Total blog Posts")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "successfully deleted",
	})

}

// Internal Use (Do not Delete)
func getCurrentUser(ctx *gin.Context) (uint, error) {
	session := sessions.Default(ctx)
	currentUser := session.Get("currentUser")
	if currentUser == nil {
		return 0, errors.New("unauthenticated")
	}

	session.Set("currentUser", currentUser)

	return currentUser.(uint), session.Save()
}
