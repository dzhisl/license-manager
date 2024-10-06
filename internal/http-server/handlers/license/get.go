package license

import (
	"fmt"

	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/dzhisl/license-manager/internal/storage/sqlite"
	"github.com/gin-gonic/gin"
)

type licenseGetterByUserId interface {
	GetLicenseById(Id string) (*sqlite.UserLicense, error)
	GetLicenseByLicense(License string) (*sqlite.UserLicense, error)
}

// GetLicenseHandler responds with license information based on either UserId or License
func GetLicenseHandler(c *gin.Context, licenseGetter licenseGetterByUserId) {
	userId := c.Query("UserId")
	license := c.Query("License")

	// Check if both parameters are empty
	if userId == "" && license == "" {
		response.InvalidInputError(c, fmt.Errorf("either UserId or License parameter is required"))
		return
	}

	var licenseData *sqlite.UserLicense
	var err error

	// Check for UserId first, then License
	if userId != "" {
		licenseData, err = licenseGetter.GetLicenseById(userId)
	} else if license != "" {
		licenseData, err = licenseGetter.GetLicenseByLicense(license)
	}

	if err != nil {
		response.InternalError(c, "Failed to get license", err)
		return
	}

	response.Ok(c, "License received", licenseData)
}
