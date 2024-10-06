package license

import (
	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/dzhisl/license-manager/internal/storage/sqlite"
	"github.com/gin-gonic/gin"
)

type allLicensesGetter interface {
	GetAllLicenses() ([]sqlite.UserLicense, error)
}

func GetAllLicensesHandler(c *gin.Context, allLicensesGetter allLicensesGetter) {
	allLicenses, err := allLicensesGetter.GetAllLicenses()
	if err != nil {
		response.InternalError(c, "Failed to get licenses", err)
		return
	}
	response.Ok(c, "success", allLicenses)
}
