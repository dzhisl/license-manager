package server

import (
	"log"

	"golang.org/x/exp/slog"

	"os"

	"github.com/dzhisl/license-manager/internal/config"
	"github.com/dzhisl/license-manager/internal/http-server/handlers/license"
	"github.com/dzhisl/license-manager/internal/http-server/handlers/ping"
	"github.com/dzhisl/license-manager/internal/http-server/middleware"
	"github.com/dzhisl/license-manager/internal/storage/sqlite"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the Gin router
func SetupRouter(storage *sqlite.Storage, AuthData *config.AuthData, sllogger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	setupGinLogs()
	r := gin.Default()
	r.Use(middleware.RequestLogger(sllogger))
	registerPublicRoutes(r, storage)

	// Using the API key for authentication
	protected := r.Group("/")
	protected.Use(middleware.APIKeyAuthMiddleware(AuthData.ApiKey)) // Use API key middleware

	registerProtectedRoutes(protected, storage)

	return r
}

// setupGinLogs set up logs for gin to be logged into logs/logs.log
func setupGinLogs() {
	f, err := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	gin.DefaultWriter = f
	gin.DefaultErrorWriter = f
}

// registerPublicRoutes registers the routes that do not require authentication.
func registerPublicRoutes(r *gin.Engine, storage *sqlite.Storage) {
	r.GET("/ping", ping.PingHandler)
	r.POST("/bind-license", func(c *gin.Context) { license.BindLicenseHandler(c, storage) })
	r.POST("/unbind-license", func(c *gin.Context) { license.UnbindLicenseHandler(c, storage) })
	r.POST("/validate-license", func(c *gin.Context) { license.ValidateLicenseHandler(c, storage) })
}

// registerProtectedRoutes registers the routes that require authentication.
func registerProtectedRoutes(authorized *gin.RouterGroup, storage *sqlite.Storage) {
	authorized.GET("/get", func(c *gin.Context) { license.GetLicenseHandler(c, storage) })
	authorized.GET("/all-licenses", func(c *gin.Context) { license.GetAllLicensesHandler(c, storage) })
	authorized.POST("/add-license", func(c *gin.Context) { license.AddLicenseHandler(c, storage) })
	authorized.POST("/del-license", func(c *gin.Context) { license.DeletelicenseHandler(c, storage) })
	authorized.POST("/freeze-license", func(c *gin.Context) { license.FreezeLicenseHandler(c, storage) })
	authorized.POST("/unfreeze-license", func(c *gin.Context) { license.UnfreezeLicenseHandler(c, storage) })
	authorized.POST("/renew-license", func(c *gin.Context) { license.RenewLicenseHandler(c, storage) })
}
