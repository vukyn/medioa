package xhttp

import (
	"github.com/gin-gonic/gin"
)

func Ok(ctx *gin.Context, result any) {
	ctx.JSON(STATUS_OK, gin.H{
		"success": true,
		"data":    result,
		"status":  Text(STATUS_OK),
	})
}

func Created(ctx *gin.Context, result any) {
	ctx.JSON(STATUS_OK, gin.H{
		"success": true,
		"data":    result,
		"status":  Text(STATUS_CREATED),
	})
}

func BadRequest(ctx *gin.Context, err error) {
	ctx.JSON(STATUS_OK, gin.H{
		"error": gin.H{
			"code":    STATUS_BAD_REQUEST,
			"message": err.Error(),
			"status":  Text(STATUS_BAD_REQUEST),
		},
	})
}

func Internal(ctx *gin.Context, err error) {
	ctx.JSON(STATUS_OK, gin.H{
		"error": gin.H{
			"code":    STATUS_INTERNAL_SERVER_ERROR,
			"message": err.Error(),
			"status":  Text(STATUS_INTERNAL_SERVER_ERROR),
		},
	})
}
