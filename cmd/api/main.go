package main

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luponetn/.enx/internal/config"
	"github.com/luponetn/.enx/internal/db"
	"github.com/luponetn/.enx/internal/logger"
	"github.com/luponetn/.enx/internal/utils"
)

type App struct {
	DBPool *pgxpool.Pool
}

func main() {

	//setup and load app config
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("could not load the app config", "error", err)
		return
	}

	//setup logger
	env := utils.ExtractKeyFromEnv("APP_ENV", "development")
	if env == "" {
		slog.Error("could not retrieve the app environment for the logger initialization")
	}
	logger.InitLogger(env)

	//set up db and connect
	dbPool, err := db.ConnectDB(cfg.DbUrl)
	if err != nil {
		slog.Error("could not connect to the database", "error", err)
		return
	}

	//setup app struct
	_ = &App{
		DBPool: dbPool,
	}

	//create router and startup server
	router := CreateRouter()
	if router == nil {
		slog.Error("could not create the router")
		return
	}

	slog.Info("starting the server", "port", cfg.Port)

	if err := StartServer(router, cfg.Port); err != nil {
		slog.Error("could not start the server", "error", err)
		return
	}

}
