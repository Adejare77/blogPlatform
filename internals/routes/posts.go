package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var PostRoutes = func(r *gin.RouterGroup) {
	r.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.MyPosts)
	r.GET("/posts/drafts", controllers.MyDrafts)
	r.GET("/posts/drafts/:id", controllers.GetDraft)
	r.PUT("/posts/:id", controllers.UpdatePost)
	r.PUT("/posts/drafts/:id", controllers.UpdatePost)
	r.DELETE("/posts/:id", controllers.DeletePost)
	r.GET("/logout", controllers.Logout)
}
