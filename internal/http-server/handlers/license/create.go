package license

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/dzhisl/license-manager/internal/http-server/response"
	"github.com/gin-gonic/gin"
)

// licenseAdder defines an interface for adding a license
type licenseAdder interface {
	AddLicense(license, username, status string, hwid *string, expiresAt time.Time) (int64, error)
}

// AddInputData represents the incoming data structure
type AddInputData struct {
	UserId string `json:"user_id" binding:"required"`
}

// OutputData represents the data structure to return
type OutputData struct {
	UserId    string    `json:"user_id"`
	License   string    `json:"license"`
	Status    string    `json:"status"`
	Hwid      string    `json:"hwid"`
	ExpiresAt time.Time `json:"expires_at"`
}

// seededRand is used for generating random licenses.
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const licenseLength = 10

// AddLicenseHandler generates a random license, sets the status to "active",
// sets the HWID to an empty string, and sets expires_at to 1 month from now.
// It also handles errors in both input validation and license addition.
func AddLicenseHandler(c *gin.Context, licenseAdder licenseAdder) {
	var input AddInputData

	// Bind incoming JSON data to the AddInputData struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Generate a random license and prepare default values
	license := generateRandomLicense(licenseLength)
	status := "active"
	hwid := ""
	expiresAt := time.Now().AddDate(0, 1, 0)

	// Attempt to add the license to the system via the licenseAdder
	if _, err := licenseAdder.AddLicense(license, input.UserId, status, &hwid, expiresAt); err != nil {
		response.InternalError(c, "Failed to add license", err)
		return
	}

	// Prepare response data
	output := OutputData{
		UserId:    input.UserId,
		License:   license,
		Status:    status,
		Hwid:      hwid,
		ExpiresAt: expiresAt,
	}

	// Respond with success and generated data
	response.Ok(c, "License added!", output)
}

// generateRandomLicense generates a random alphanumeric license of the given length
func generateRandomLicense(length int) string {
	b := make([]byte, length)
	charsetLength := len(charset) // Cache length outside of the loop for efficiency
	for i := range b {
		b[i] = charset[seededRand.Intn(charsetLength)]
	}
	return string(b)
}
