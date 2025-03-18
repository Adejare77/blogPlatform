package routes

import (
	"github.com/Adejare77/blogPlatform/internals/controllers"
	"github.com/gin-gonic/gin"
)

var UserRoutes = func(r *gin.RouterGroup) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/index", controllers.AllPosts)
	r.GET("/posts/:id", controllers.GetPost)
}
