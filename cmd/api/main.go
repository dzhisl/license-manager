package main

import (
	"os"

	"golang.org/x/exp/slog" // Change this

	"github.com/dzhisl/license-manager/internal/config"
	"github.com/dzhisl/license-manager/internal/http-server/server"
	"github.com/dzhisl/license-manager/internal/lib/logger"
	"github.com/dzhisl/license-manager/internal/lib/logger/sl"
	"github.com/dzhisl/license-manager/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()
	logger := logger.SetupLogger()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	r := server.SetupRouter(storage, &cfg.AuthData, logger)
	logger.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))
	if err := r.Run(cfg.HTTPServer.Address); err != nil {
		logger.Error("server failed", sl.Err(err))
		os.Exit(1)
	}
}
