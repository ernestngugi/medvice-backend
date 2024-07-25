package webutils

import (
	"github.com/ernestngugi/medvice-backend/internal/apperror"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, appError *apperror.Error) {

	jsonResponse := map[string]any{
		"status":        false,
		"error_message": appError.Error(),
	}

	c.JSON(appError.HttpStatusCode(), jsonResponse)
}
