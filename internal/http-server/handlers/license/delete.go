package license

import (
	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

// licenseAdder defines an interface for adding a license
type licenseDeleter interface {
	DeleteLicenseById(username string) error
}

// DeleteInputData represents the incoming data structure
type DeleteInputData struct {
	UserId string `json:"user_id" binding:"required"`
}

func DeletelicenseHandler(c *gin.Context, licenseDeleter licenseDeleter) {
	var input DeleteInputData

	if err := c.ShouldBindJSON(&input); err != nil {
		response.InvalidInputError(c, err)
		return
	}

	err := licenseDeleter.DeleteLicenseById(input.UserId)
	if err != nil {
		response.InternalError(c, "Failed to delete license", err)
		return
	}

	output := map[string]string{
		"UserId": input.UserId,
	}

	response.Ok(c, "License deleted successfully", output)
}
