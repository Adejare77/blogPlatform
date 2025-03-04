package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adejare77/blogPlatform/config"
	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/Adejare77/blogPlatform/internals/models"
	"github.com/Adejare77/blogPlatform/internals/schemas"
	"github.com/Adejare77/blogPlatform/internals/utilities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Register(ctx *gin.Context) {
	var user schemas.User

	if err := ctx.ShouldBind(&user); err != nil {
		utilities.Validator(ctx, err)
		return
	}

	// Create User
	if err := models.CreateUser(&user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			handlers.BadRequest(ctx, "email already exists", err)
			return
		}
		handlers.InternalServerError(ctx, err)
		return
	}

	// log registered User
	logrus.WithFields(logrus.Fields{
		"email": user.Email,
	}).Info("User registered successfully")

	ctx.JSON(http.StatusCreated, gin.H{
		"email":    user.Email,
		"fullName": user.Name,
		"status":   "Successfuly Registered",
	})
}

func Login(ctx *gin.Context) {
	type login struct {
		Email    string `binding:"required,email" json:"email"`
		Password string `binding:"required" json:"password"`
	}

	var user login
	if err := ctx.ShouldBind(&user); err != nil {
		utilities.Validator(ctx, err)
		return
	}

	userInfo, ok := models.GetUserInfo(user.Email)
	if !ok {
		handlers.Unauthorized(ctx, "Incorrect email or password", userInfo)
		return
	}
	if err := utilities.ComparePassword(user.Password, userInfo.Password); err != nil {
		handlers.Unauthorized(ctx, "Incorrect email or password", userInfo)
		return
	}

	config.CreateSession(ctx, userInfo.ID)

	ctx.JSON(http.StatusOK, "Successfully Login")
}

func Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("currentUser")
	fmt.Println("**********************")
	fmt.Println(userID)
	fmt.Println("**********************")

	if !exists {
		handlers.Unauthorized(ctx, "Unauthorized", "User Session Not Found")
		return
	}

	config.DeleteSession(ctx, userID.(string))

	ctx.Redirect(http.StatusSeeOther, "index")
}
