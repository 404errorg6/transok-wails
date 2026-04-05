package middleware

import (
	"transok/backend/domain/resp"
	"transok/backend/services"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler is the captcha verification middleware
func CaptchaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaInStore := services.Share().GetCaptcha()
		/* No captcha set, skip directly */
		if captchaInStore == "" {
			c.Next()
			return
		}

		// Get captcha from request header
		captchaInHeader := c.GetHeader("Captcha-Key")
		if captchaInHeader == "" {
			resp.Forbidden().WithMsg("Captcha is required").Out()
			c.Abort()
			return
		}

		if captchaInHeader != captchaInStore {
			resp.Forbidden().WithMsg("Captcha is incorrect").Out()
			c.Abort()
			return
		}

		c.Next()
	}
}
