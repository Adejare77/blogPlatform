package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var PostRoutes = func(r *gin.RouterGroup) {
	// Published Blogs
	r.PATCH("/posts/:post_id", controllers.UpdatePost)
	r.DELETE("/posts/:post_id", controllers.DeletePost)
	r.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.GetUserPosts)

	r.GET("/logout", controllers.Logout)
}
