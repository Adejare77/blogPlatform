package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		"status": statusCode,
		"error":  errorMessage,
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

func Forbidden(ctx *gin.Context, msg string, details any) {
	handleError(
		ctx,
		http.StatusForbidden,
		"FORBIDDEN",
		msg,
		details,
	)
}

func Validator(ctx *gin.Context, err error) {
	var errorDetails []string
	for _, fieldError := range err.(validator.ValidationErrors) {
		validationError := "validation failed: "
		// source := errorSource(fieldError)
		if fieldError.Tag() == "required" {
			errorDetails = append(
				errorDetails,
				fmt.Sprintf(
					validationError+"missing `%s` field on `%s`",
					fieldError.Field(),
					fieldError.Tag(),
				))
		} else if fieldError.Tag() == "oneof" {
			errorDetails = append(errorDetails, validationError+"`status` parameter can only be `draft` or `published`")
		} else if fieldError.Tag() == "uuid" {
			errorDetails = append(errorDetails, validationError+"invalid post uuid")
		} else if fieldError.Tag() == "status" {
			errorDetails = append(errorDetails, validationError+"missing status query")
		} else if fieldError.Tag() == "email" {
			errorDetails = append(errorDetails, validationError+"invalid email")
		} else {
			errorDetails = append(errorDetails, validationError+fieldError.Error())
		}
	}
	handleError(
		ctx,
		http.StatusBadRequest,
		"Validation Error",
		errorDetails,
		err,
	)
}
