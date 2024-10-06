package license

import (
	"net/http"
	"time"

	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/dzhisl/license-manager/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

// licenseValidator defines the required methods for validating a license
type licenseValidator interface {
	GetLicenseByLicense(license string) (*sqlite.UserLicense, error)
	BindHwidToLicenseByLicense(license, hwid string) error
}

// validateInputData represents the incoming data for validation
type validateInputData struct {
	License string `json:"license" binding:"required"`
	HWID    string `json:"hwid" binding:"required"`
}

// ValidateLicenseHandler handles license validation requests
func ValidateLicenseHandler(c *gin.Context, licenseValidator licenseValidator) {
	var input validateInputData

	// Bind and validate input data
	if err := c.ShouldBindJSON(&input); err != nil {
		response.InvalidInputError(c, err)
		return
	}

	// Retrieve license information by license ID
	licenseData, err := licenseValidator.GetLicenseByLicense(input.License)
	if err != nil {
		response.InternalError(c, "failed to get license", err)
		return
	}

	// Validate if the license is active
	if licenseData.Status != "active" {
		response.Error(c, "license is not active", http.StatusForbidden, nil)
		return
	}

	// Validate if the license has expired
	if time.Now().After(licenseData.ExpiresAt) {
		response.Error(c, "license has expired", http.StatusForbidden, nil)
		return
	}

	// Validate HWID and bind if necessary
	if licenseData.HWID == nil || *licenseData.HWID == "" {
		// HWID is empty, bind it to the license
		if err := licenseValidator.BindHwidToLicenseByLicense(input.License, input.HWID); err != nil {
			response.InternalError(c, "failed to bind HWID to license", err)
			return
		}
	} else if *licenseData.HWID != input.HWID {
		// HWID mismatch
		response.Error(c, "HWID does not match", http.StatusForbidden, nil)
		return
	}

	// License is valid, respond with success
	response.Ok(c, "license is valid!", licenseData)
}
