package main

import (
	"log/slog"

	"github.com/luponetn/.enx/internal/config"
	"github.com/luponetn/.enx/internal/db"
	"github.com/luponetn/.enx/internal/logger"
	"github.com/luponetn/.enx/internal/utils"
)

func main() {

	// setup and load app config
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("could not load the app config", "error", err)
		return
	}

	// setup logger
	env := utils.ExtractKeyFromEnv("APP_ENV", "development")
	logger.InitLogger(env)

	// set up db and connect
	dbPool, err := db.ConnectDB(cfg.DbUrl)
	if err != nil {
		slog.Error("could not connect to the database", "error", err)
		return
	}

	// setup queries and server
	queries := db.New(dbPool)
	server := NewServer(queries, cfg)

	// create router and startup server
	router := CreateRouter(server)

	slog.Info("starting the server", "port", cfg.Port)

	if err := StartServer(router, cfg.Port); err != nil {
		slog.Error("could not start the server", "error", err)
		return
	}
}