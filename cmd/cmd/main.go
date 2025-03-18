package main

import (
	"log"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/middlesware"
	"github.com/Adejare77/blogPlatform/internals/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration (this sets config.DB, etc.)
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to Initialize Server: %v", err)
	}

	app := gin.Default()
	app.Use(sessions.Sessions("blogPost", config.SessionStore))

	userRoutes := app.Group("/")
	postRoutes := app.Group("/", middlesware.AuthMiddleware())
	likeRoutes := app.Group("/", middlesware.AuthMiddleware())

	routes.UserRoutes(userRoutes)
	routes.PostRoutes(postRoutes)
	routes.LikesRoutes(likeRoutes)

	log.Print("Server running successfully")
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Could not Run on Port 3000")
	}

}
