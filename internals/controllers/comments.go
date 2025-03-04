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
	userID := ctx.MustGet("currentUser").(uint)
	postID := ctx.Param("id")
	parentID := ctx.Param("comment_id")

	var comment schemas.Comment

	if err := ctx.BindJSON(&comment); err != nil {
		handlers.BadRequest(ctx, err.Error(), err)
		return
	}

	// Add PostID and UserID
	comment.PostID = postID
	comment.UserID = userID

	if parentID != "" { // its a reply to a comment
		comment.ParentID = &parentID
	}

	if err := models.PostComment(comment); err != nil {
		handlers.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Comment Added Successfully")
}

func GetComments(ctx *gin.Context) {
	postID := ctx.Param("id")

	allComments, err := models.GetCommentsByPostID(postID)
	if err != nil {
		handlers.InternalServerError(ctx, err.Error())
	}

	ctx.JSON(http.StatusOK, allComments)
}

func UpdateComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	commentID := ctx.Param("comment_id")

	var content map[string]string
	if err := ctx.ShouldBind(&content); err != nil {
		handlers.BadRequest(ctx, err.Error(), err)
		return
	}

	if content["content"] == "" {
		handlers.BadRequest(ctx, "`content` field missing", "`Content` field is not Provided")
		return
	}
	err := models.UpdateComment(userID, commentID, content)
	if err == nil {
		ctx.Redirect(http.StatusSeeOther, commentID)
	} else if strings.Contains(err.Error(), "not found") {
		ctx.JSON(http.StatusNotFound, "Record not found")
	} else {
		handlers.InternalServerError(ctx, err)
	}
}

func DeleteComment(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	commentID := ctx.Param("comment_id")

	err := models.DeleteComment(userID, commentID)
	if err == nil {
		ctx.JSON(http.StatusOK, "Comment Deleted Successfully")
	} else if strings.Contains(err.Error(), "not found") {
		ctx.JSON(http.StatusNotFound, "Comment Not Found")
	} else {
		handlers.InternalServerError(ctx, err)
	}
}

// func ReplyAComment(ctx *gin.Context) {
// 	userID := ctx.MustGet("currentUser").(uint)
// 	postID := ctx.Param("id")
// 	commentID := ctx.Param("comment_id")

// 	var reply schemas.Comment
// 	if err :=

// }
