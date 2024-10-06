package license

import (
	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

type LicenseFreezer interface {
	FreezeLicenseById(UserId string) error
	UnfreezeLicenseById(UserId string) error
}

// FreezeInputData represents the incoming data structure
type FreezeInputData struct {
	UserId string `json:"user_id" binding:"required"`
}

// handleLicenseAction is a generic handler for freezing or unfreezing licenses
func handleLicenseAction(c *gin.Context, licenseAction func(string) error, successMessage string) {
	var input FreezeInputData
	if err := c.ShouldBindJSON(&input); err != nil {
		response.InvalidInputError(c, err)
		return
	}

	// Execute the passed license action (freeze/unfreeze)
	if err := licenseAction(input.UserId); err != nil {
		response.InternalError(c, "Operaion failed", err)
		return
	}

	// Respond with success message
	response.Ok(c, successMessage, nil)
}

// FreezeLicenseHandler handles the freezing of a license
func FreezeLicenseHandler(c *gin.Context, licenseFreezer LicenseFreezer) {
	handleLicenseAction(c, licenseFreezer.FreezeLicenseById, "License frozen successfully")
}

// UnfreezeLicenseHandler handles the unfreezing of a license
func UnfreezeLicenseHandler(c *gin.Context, licenseFreezer LicenseFreezer) {
	handleLicenseAction(c, licenseFreezer.UnfreezeLicenseById, "License unfrozen successfully")
}
