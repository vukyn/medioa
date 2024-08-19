package ratelimiter

import (
	"github.com/didip/tollbooth/v7"
	"github.com/gin-gonic/gin"
)

func LimitPerSecond(max int) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(float64(max), nil)
	return func(c *gin.Context) {
		if err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request); err != nil {
			c.Data(err.StatusCode, lmt.GetMessageContentType(), []byte(err.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
