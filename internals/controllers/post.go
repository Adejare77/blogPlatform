package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"github.com/Adejare77/blogPlatform/internals/utilities"
	"github.com/gin-gonic/gin"
)

func CreatePost(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)

	var post schemas.Post

	if err := ctx.ShouldBind(&post); err != nil {
		fmt.Println(err)
		utilities.Validator(ctx, err)
		return
	}

	post.UserID = userID

	if err := models.CreatePost(&post); err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"title":   post.Title,
		"content": post.Content,
		"status":  post.Status,
	})
}

func AllPosts(ctx *gin.Context) {
	posts, err := models.GetAllPosts()
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func MyPosts(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)

	posts, err := models.GetPosts(userID)
	if err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	if len(posts) == 0 {
		ctx.JSON(http.StatusNotFound, "No Record Found")
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func GetPost(ctx *gin.Context) {
	ctx.MustGet("currentUser")
	postID := ctx.Param("id")

	post, err := models.GetPost(postID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusOK, "Record Not Found")
			return
		}
		utilities.Validator(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func UpdatePost(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	postID := ctx.Param("id")

	var data map[string]interface{}
	if err := ctx.ShouldBind(&data); err != nil {
		utilities.Validator(ctx, err)
		return
	}

	if err := models.UpdatePost(userID, postID, data); err != nil {
		handlers.InternalServerError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, postID)
}

func DeletePost(ctx *gin.Context) {
	userID := ctx.MustGet("currentUser").(uint)
	postID := ctx.Param("id")

	status_code, err := models.DeletePost(userID, postID)

	if err == nil {
		ctx.JSON(status_code, "Successfully Deleted")
		return
	}

	ctx.JSON(status_code, err.Error())
}
