package controllers

import (
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/gin-gonic/gin"
)

func LikedPosts(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)

	posts, err := models.GetLikedPosts(userID)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}
	if len(posts) == 0 {
		ctx.JSON(http.StatusNotFound, "Record Not Found")
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func LikePostOrComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	postID := ctx.Param("id")
	commentID := ctx.Param("comment_id")

	user, err := models.GetUserByPostID(postID)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}
	if user.ID == userID {
		handlers.BadRequest(ctx, "Author Cannot like their own post", err)
		return
	}

	postOrComment := "Post"
	postIDOrCommentID := postID

	if commentID != "" {
		postOrComment = "Comment"
		postIDOrCommentID = commentID
	}

	err = models.LikePostOrComment(userID, postIDOrCommentID, postOrComment)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Post successfully Liked")
}

func UnlikedPostOrComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	postID := ctx.Param("id")
	commentID := ctx.Param("comment_id")

	postOrComment := "Post"
	postIDOrCommentID := postID

	if commentID != "" {
		postOrComment = "Comment"
		postIDOrCommentID = commentID
	}

	err := models.UnlikePostOrComment(userID, postIDOrCommentID, postOrComment)
	if err == nil {
		ctx.JSON(http.StatusOK, "Post Successfully Unliked")
	} else if strings.Contains(err.Error(), "not found") {
		ctx.JSON(http.StatusNotFound, err.Error())
	} else {
		handlers.InternalServerError(ctx, err)
	}

}
