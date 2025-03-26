package main

import (
	"log"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/middlesware"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration (this sets config.DB, etc.)
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to Initialize Server: %v", err)
	}

	totalPosts, err := models.TotalPosts()
	if err != nil {
		log.Fatalf("Faild to fetch Total Posts")
	}

	config.TotalPosts = totalPosts

	app := gin.Default()
	app.Use(sessions.Sessions("blogPost", config.SessionStore))

	postRoutes := app.Group("/", middlesware.AuthMiddleware())
	likeRoutes := app.Group("/", middlesware.AuthMiddleware())
	commentRoutes := app.Group("/", middlesware.AuthMiddleware())
	userRoutes := app.Group("/")

	routes.LikesRoutes(likeRoutes)
	routes.PostRoutes(postRoutes)
	routes.UserRoutes(userRoutes)
	routes.CommentRoutes(commentRoutes)

	log.Print("Server running successfully")
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Could not Run on Port 3000")
	}

}
