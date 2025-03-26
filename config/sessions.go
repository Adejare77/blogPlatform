package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type sessionConfig struct {
	Size      int
	Address   string
	Password  string
	secretKey string
	MaxAge    int
}

var SessionStore redis.Store

func ConnectSession() error {
	cfg := loadSessionConfig()
	client, err := redis.NewStore(
		cfg.Size,
		"tcp",
		cfg.Address,
		cfg.Password,
		[]byte(cfg.secretKey),
	)

	if err != nil {
		return fmt.Errorf("(Session Initialization Error) %v", err)
	}

	// global session options
	client.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.MaxAge,
		HttpOnly: true,  // Prevent client-side JS access
		Secure:   false, // Set to true if using HTTPS
	})

	SessionStore = client

	return nil
}

func CreateSession(ctx *gin.Context, userID uint) {
	session := sessions.Default(ctx)

	session.Set("currentUser", userID)

	if err := session.Save(); err != nil {
		handlers.InternalServerError(ctx, "Error Saving User session")
		return
	}
}

func DeleteSession(ctx *gin.Context, userID uint) {
	session := sessions.Default(ctx)

	session.Clear()
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})

	if err := session.Save(); err != nil {
		handlers.InternalServerError(ctx, "Error Saving Deleted User session")
	}
}

func loadSessionConfig() *sessionConfig {
	size := os.Getenv("REDIS_SIZE")
	addr := os.Getenv("REDIS_ADDRESS")
	pwd := os.Getenv("REDIS_PASSWORD")
	secretKey := os.Getenv("REDIS_SECRETKEY")
	maxAge := os.Getenv("REDIS_MAX_AGE")

	maxIdleConn, err := strconv.Atoi(size)
	if err != nil {
		fmt.Println("Warning: Incorrect size value; defaults to 10")
		maxIdleConn = 10
	}

	age, err := strconv.Atoi(maxAge)
	if err != nil {
		fmt.Println("Warning: Incorrect MAX_AGE value; defaults to 600s")
		maxIdleConn = 10
		age = 600
	}

	return &sessionConfig{
		Size:      maxIdleConn,
		Address:   addr,
		Password:  pwd,
		secretKey: secretKey,
		MaxAge:    age,
	}
}
