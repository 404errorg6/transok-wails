package middleware

import (
	"runtime/debug"
	"transok/backend/domain/resp"

	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {

			// Check if the panic is a predefined Error response body
			if e, ok := err.(resp.Err); ok {
				// If so, return directly
				if !e.Success {
					debug.PrintStack()
					c.JSON(400, e)
					c.Abort()
					return
				}
				c.JSON(200, e)
			} else {
				debug.PrintStack()
				c.JSON(500, resp.ServerErr())
			}
			c.Abort()
		}
	}()

	c.Next()
}
