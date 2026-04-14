package middleware

import (
	reflect "reflect"
	"strings"
	"transok/backend/domain/resp"

	"github.com/gin-gonic/gin"
)

func Valid(dto interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new object using the passed-in DTO
		data := reflect.New(reflect.TypeOf(dto)).Interface()

		if err := c.ShouldBindJSON(data); err != nil {
			resp.DataFormatErr().WithData(strings.Split(err.Error(), "\n")).Out()
			c.Abort()
			return
		}
		c.Set("dto", data)
		c.Next()
	}
}
