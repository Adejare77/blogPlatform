package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Details    interface{}
}

func handleError(ctx *gin.Context, statusCode int, errorCode string, errorMessage interface{}, errorDetails interface{}) {
	// For Developers
	logrus.WithFields(logrus.Fields{
		"statusCode":   statusCode,
		"errorCode":    errorCode,
		"errorDetails": errorDetails,
	})

	// For users/clients
	ctx.JSON(statusCode, gin.H{
		"error_code": errorCode,
		"msg":        errorMessage,
	})
}

func BadRequest(ctx *gin.Context, msg interface{}, details interface{}) {
	handleError(
		ctx,
		http.StatusBadRequest,
		"BAD_REQUEST",
		msg,
		details,
	)
}

func InternalServerError(ctx *gin.Context, details interface{}) {
	handleError(
		ctx,
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"An internal Server Error Occurred. Try Again later",
		details,
	)
}

func Unauthorized(ctx *gin.Context, msg string, details interface{}) {
	handleError(
		ctx,
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		msg,
		details,
	)
}
