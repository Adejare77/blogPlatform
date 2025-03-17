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
			fmt.Println(fieldError.Tag())
			if fieldError.Tag() == "required" {
				errorDetails = append(errorDetails, fmt.Sprintf("The `%s` field is required", fieldError.Field()))
			} else if fieldError.Tag() == "oneof" {
				errorDetails = append(errorDetails, "The `status` field can only be `unpublished` or `published`")
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
