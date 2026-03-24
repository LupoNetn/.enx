package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/luponetn/enx/internal/utils"
)

type Config struct {
	DbUrl            string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func LoadConfig() (*Config, error) {
  godotenv.Load()
  
  DbUrl := utils.ExtractKeyFromEnv("DATABASE_URL", "")
  if DbUrl == "" {
	slog.Error("could not retrieve the database url from the environment variables")
	return nil, os.ErrNotExist
  }

  port := utils.ExtractKeyFromEnv("PORT", "5050")
  if port == "" {
	slog.Error("could not retrieve the port from the environment variables")
	return nil, os.ErrNotExist
  }

  jwtAccessSecret := utils.ExtractKeyFromEnv("JWT_ACCESS_SECRET", "")
  if jwtAccessSecret == "" {
	slog.Error("could not retrieve the jwt access secret from the environment variables")
	return nil, os.ErrNotExist
  }

  jwtRefreshSecret := utils.ExtractKeyFromEnv("JWT_REFRESH_SECRET", "")
  if jwtRefreshSecret == "" {
	slog.Error("could not retrieve the jwt refresh secret from the environment variables")
	return nil, os.ErrNotExist
  }

  return &Config{
	DbUrl:            DbUrl,
	Port:             port,
	JWTAccessSecret:  jwtAccessSecret,
	JWTRefreshSecret: jwtRefreshSecret,
	}, nil

}

