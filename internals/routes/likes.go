package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var LikesRoutes = func(r *gin.RouterGroup) {
	r.GET("/posts/like", controllers.LikedPosts)
	r.POST("/posts/:id/like", controllers.LikePostOrComment)
	r.DELETE("posts/:id/unlike", controllers.UnlikedPostOrComment)
	r.POST("/posts/:id/comments/:comment_id/like", controllers.LikePostOrComment)
	r.DELETE("posts/:id/comments/:comment_id/unlike", controllers.UnlikedPostOrComment)
}
