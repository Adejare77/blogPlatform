package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var CommentRoutes = func(r *gin.RouterGroup) {
	r.GET("/posts/:post_id/comments", controllers.GetPostComments)
	r.POST("/posts/:post_id/comments", controllers.CreateComment)
	r.POST("/posts/:post_id/comments/:parent_id/replies", controllers.CreateComment)
	r.PATCH("/posts/:post_id/comments/:comment_id", controllers.UpdateComment)
}
