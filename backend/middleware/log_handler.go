package middleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"transok/backend/domain/resp"

	"github.com/gin-gonic/gin"
)

func LogHandler(basePath string) gin.HandlerFunc {
	logPath := filepath.Join(basePath, "log")

	// Create log directory
	if err := os.MkdirAll(logPath, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
	}

	logFile, err := os.OpenFile(filepath.Join(basePath, "logs", "gin.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
	}
	logger := log.New(logFile, "[API-LOG]", log.LstdFlags)
	return func(c *gin.Context) {
		// Request start time
		startTime := time.Now()
		var err interface{}
		defer func() {
			if err != nil {
				errRes := resp.UnknownErr().Msg
				businessCode := resp.UnknownErr().Code
				if e, ok := err.(resp.Err); ok {
					errRes = e.Msg
					businessCode = e.Code
				}
				// Request end time
				endTime := time.Now()
				// Calculate request processing time
				processingTime := endTime.Sub(startTime)
				// Request information
				reqMethod := c.Request.Method
				reqURL := c.Request.RequestURI
				statusCode := c.Writer.Status()
				clientIp := c.ClientIP()
				logMsg := fmt.Sprintf(" <IP: ==> %s> [%s] URL: %s => [code:%d | BusinessCode:%d => \"%s\"] lasted:%v", clientIp, reqMethod, reqURL, statusCode, businessCode, errRes, processingTime)
				logger.Println(logMsg)
				panic(err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				err = r
			}
		}()

		c.Next()

	}
}

type LogType string

const (
	API LogType = "api"
)

func (l LogType) String() string {
	// Declare a LogType variable
	var api LogType = "apassi"
	return string(api)
}
