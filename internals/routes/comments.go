package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var CommentRoutes = func(r *gin.RouterGroup) {
	r.POST("/posts/:id/comments", controllers.CreateComment)
	r.GET("/posts/:id/comments", controllers.GetComments)
	r.POST("/posts/:id/comments/:comment_id/reply", controllers.CreateComment)
	r.PUT("/posts/:id/comments/:comment_id")
}
