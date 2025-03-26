package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"github.com/gin-gonic/gin"
)

func CreateLike(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(uint)
	var post schemas.PostURIParams
	var targetType schemas.LikedQueryParams

	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&targetType); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	var authorID *uint
	var err error

	if targetType.Type == "post" {
		authorID, err = models.FindPostAuthorID(post.PostID)
	} else {
		authorID, err = models.FindCommentAuthorID(post.PostID)
	}

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "Record not Found", err)
			return
		}
		handlers.InternalServerError(ctx, err)
		return
	}

	if *authorID == currentUser {
		handlers.BadRequest(ctx,
			fmt.Sprintf("you can't like your %s", targetType.Type),
			fmt.Sprintf("author tried to like their %s", targetType.Type),
		)
		return
	}

	if err := models.CreateLike(schemas.Like{
		UserID:       currentUser,
		LikeableID:   post.PostID,
		LikeableType: targetType.Type,
	}); err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "liked successfully")
}

func GetUserLikes(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	var targetType schemas.LikedQueryParams

	if err := ctx.ShouldBindQuery(&targetType); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	posts, err := models.FindLikesByUser(userID, targetType.Type)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}
	if len(posts) == 0 {
		ctx.JSON(http.StatusOK, []map[string]any{})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func DeleteLike(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	var post schemas.PostURIParams

	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	err := models.DeleteLike(userID, post.PostID)
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "deleted successfully",
	})

}
