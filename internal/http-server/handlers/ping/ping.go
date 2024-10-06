package ping

import (
	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

// PingHandler responds with "pong" when /ping is accessed
func PingHandler(c *gin.Context) {
	response.Ok(c, "pong", nil)
}
