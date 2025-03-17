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
			handlers.Unauthorized(ctx, "Login Required", "Access Unauthorized routes")
			ctx.Abort()
			return
		}

		userID, ok := currentUser.(uint)
		if !ok {
			handlers.Unauthorized(ctx, "Invalid Session", "Session data corrupted")
			ctx.Abort()
			return
		}

		// Re-save to ensure roll-over to prevent TTL while in use
		ctx.Set("currentUser", userID)

		session.Options(sessions.Options{
			MaxAge:   600,
			HttpOnly: true,
			Secure:   false,
		})

		if err := session.Save(); err != nil {
			handlers.InternalServerError(ctx, "Error saving a session")
			return
		}
	}
}
