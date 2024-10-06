package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/exp/slog"

	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// APIKeyAuthMiddleware checks for a valid API key in the request headers.
func APIKeyAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key != apiKey {
			response.Error(c, "Unauthorized", 401, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestLogger logs incoming requests and their responses using Gin's logger.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqUUID := uuid.New().String()

		reqIP := c.ClientIP()

		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var requestBody interface{}
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			requestBody = string(bodyBytes)
		}

		requestBodyJSON, err := json.Marshal(requestBody)
		if err != nil {
			requestBodyJSON = []byte(fmt.Sprintf(`"error marshalling body: %v"`, err))
		}

		c.Next()

		statusCode := c.Writer.Status()
		logMessage := fmt.Sprintf(
			"Request UUID: %s | IP: %s | Method: %s | Path: %s | Status: %d | Body: %s",
			reqUUID, reqIP, c.Request.Method, c.Request.URL.Path, statusCode, string(requestBodyJSON),
		)
		logger.Info(logMessage)
		fmt.Fprintln(gin.DefaultWriter, logMessage)
	}
}
