package controllers

import (
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"github.com/gin-gonic/gin"
)

func CreateComment(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(uint)
	var post schemas.PostURIParams
	var content schemas.CommentBody

	// Bind URI parameters (post_id and parent_id if present)
	if err := ctx.ShouldBindUri(&post); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	// Bind JSON data (content)
	if err := ctx.ShouldBind(&content); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	comment := schemas.Comment{
		AuthorID: currentUser,
		PostID:   post.PostID,
		Content:  content.Content,
	}

	if err := models.CreateComment(comment); err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "comment added",
	})
}

func GetPostComments(ctx *gin.Context) {
	type Post struct {
		PostID string `uri:"post_id" binding:"required,uuid"`
	}
	var comment Post

	if err := ctx.ShouldBindUri(&comment); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	allComments, err := models.FindCommentsByPostID(comment.PostID)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}

	if len(allComments) == 0 {
		ctx.JSON(http.StatusOK, []map[string]any{})
		return
	}

	ctx.JSON(http.StatusOK, allComments)
}

func UpdateComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	var commentParam schemas.CommentUriParam
	var commentBody schemas.CommentBody

	if err := ctx.ShouldBindUri(&commentParam); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	if err := ctx.ShouldBindJSON(&commentBody); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	err := models.UpdateComment(userID, commentParam, commentBody)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "record not found", err)
			return
		}
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "updated successfully",
	})
}

func DeleteComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	var comment schemas.CommentUriParam

	if err := ctx.ShouldBindUri(&comment); err != nil {
		handlers.Validator(ctx, err)
		return
	}

	err := models.DeleteComment(userID, comment)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			handlers.BadRequest(ctx, "Record not Found", err)
			return
		}
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "comment deleted",
	})
}
