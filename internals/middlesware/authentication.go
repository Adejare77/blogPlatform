package middlesware

import (
	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		currentUser := session.Get("currentUser")

		if currentUser == nil {
			handlers.Unauthorized(ctx, "login required", "Access Unauthorized routes")
			ctx.Abort()
			return
		}

		userID, ok := currentUser.(uint)
		if !ok {
			handlers.Unauthorized(ctx, "Invalid Session", "Session data corrupted")
			ctx.Abort()
			return
		}

		ctx.Set("currentUser", userID)     // set key for next function
		session.Set("currentUser", userID) // roll-over key

		if err := session.Save(); err != nil {
			handlers.InternalServerError(ctx, "error saving a session")
			return
		}

		ctx.Next()
	}
}
