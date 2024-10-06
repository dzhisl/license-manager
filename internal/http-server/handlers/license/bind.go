package license

import (
	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

// LicenseBinder defines an interface for binding a license to an HWID
type LicenseBinder interface {
	BindHwidToLicenseByLicense(license, hwid string) error
}

// LicenseUnbinder defines an interface for unbinding a license from an HWID
type LicenseUnbinder interface {
	UnbindHwidFromLicense(license string) error
}

// LicenseActionInput represents the common incoming data structure for license actions
type LicenseActionInput struct {
	License string `json:"license" binding:"required"`
	HWID    string `json:"hwid,omitempty"` // HWID is optional, only used for binding
}

// processLicenseAction handles the common logic for binding and unbinding licenses
// It takes an action function, a success message, and manages error handling.
func processLicenseAction(c *gin.Context, action func(input LicenseActionInput) error, successMessage string) bool {
	var input LicenseActionInput

	// Validate incoming JSON data
	if err := c.ShouldBindJSON(&input); err != nil {
		response.InvalidInputError(c, err)
		return false
	}

	// Perform the provided license action
	if err := action(input); err != nil {
		response.InternalError(c, "Operation failed", err)
		return false
	}

	response.Ok(c, successMessage, nil)
	return true
}

// BindLicenseHandler handles binding a license to an HWID
func BindLicenseHandler(c *gin.Context, licenseBinder LicenseBinder) {
	// Define the action for binding a license to an HWID
	action := func(input LicenseActionInput) error {
		return licenseBinder.BindHwidToLicenseByLicense(input.License, input.HWID)
	}
	// Process the action
	processLicenseAction(c, action, "License bound successfully")
}

// UnbindLicenseHandler handles unbinding a license from an HWID
func UnbindLicenseHandler(c *gin.Context, licenseUnbinder LicenseUnbinder) {
	// Define the action for unbinding a license
	action := func(input LicenseActionInput) error {
		return licenseUnbinder.UnbindHwidFromLicense(input.License)
	}
	// Process the action
	processLicenseAction(c, action, "License unbound successfully")
}
