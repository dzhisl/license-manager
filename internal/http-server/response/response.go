package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 500 error response wrapper
func InternalError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}

// 200 OK response wrapper
func Ok(c *gin.Context, message string, output interface{}) {

	response := gin.H{
		"message": message,
	}

	if output != nil {
		response["data"] = output
	}

	c.JSON(http.StatusOK, response)
}

// 400 insufficient json data error wrapper
func InvalidInputError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid input data",
		"details": err.Error(),
	})
}

// custom error wrapper
func Error(c *gin.Context, errorMessage string, status int, err error) {

	response := gin.H{
		"error": errorMessage,
	}

	if err != nil {
		response["details"] = err.Error()
	}

	c.JSON(status, response)
}
