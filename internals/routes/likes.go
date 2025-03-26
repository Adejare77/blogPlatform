package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var LikesRoutes = func(r *gin.RouterGroup) {
	r.GET("/users/:user_id/likes", controllers.GetUserLikes)

	r.POST("/posts/:post_id/likes", controllers.CreateLike)
	r.DELETE("/posts/:post_id/likes", controllers.DeleteLike)

	r.POST("/comments/:comment_id/likes", controllers.CreateLike)
	r.DELETE("/comments/:comment_id/likes", controllers.DeleteLike)

}
