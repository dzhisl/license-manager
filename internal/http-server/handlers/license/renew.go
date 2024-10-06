package license

import (
	"time"

	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

type licenseRenewer interface {
	RenewLicenseById(UserId string, days int) (expirationTime time.Time, err error)
}

type RenewInputData struct {
	UserId string `json:"user_id" binding:"required"`
	Days   int    `json:"days" binding:"required"`
}

func RenewLicenseHandler(c *gin.Context, licelicenseRenewer licenseRenewer) {
	var input RenewInputData

	if err := c.ShouldBindJSON(&input); err != nil {
		response.InvalidInputError(c, err)
		return
	}

	expirationTime, err := licelicenseRenewer.RenewLicenseById(input.UserId, input.Days)
	if err != nil {
		response.InternalError(c, "Failed to renew license", err)
		return
	}

	response.Ok(c, "License renewed successfully", expirationTime)
}
