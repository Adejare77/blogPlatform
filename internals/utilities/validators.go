package utilities

import (
	"fmt"

	"github.com/Adejare77/blogPlatform/internals/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validator(ctx *gin.Context, err error) {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		var errorDetails []string
		for _, fieldError := range validationError {
			if fieldError.Tag() == "required" {
				errorDetails = append(errorDetails, fmt.Sprintf("The `%s` field is required", fieldError.Field()))
			} else if fieldError.Tag() == "status" {
				errorDetails = append(errorDetails, fmt.Sprintf("The `%s` field can only be `draft` or `published`", fieldError.Field()))
			}
		}
		handlers.BadRequest(ctx, errorDetails, errorDetails)
	} else {
		handlers.BadRequest(ctx,
			"Invalid Request",
			err.Error(),
		)
	}
}
