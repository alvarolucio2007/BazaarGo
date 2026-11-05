package main

import (
	"log"
	"os"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/logger"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/postgres"
	"github.com/alvarolucio2007/bazaargo/internal/app"
)

func main() {
	logger := logger.NewSlog()
	logger.Info("starting app")
	symmetricKey, err := app.GenerateSymmetricKey()
	if err != nil {
		logger.Error("main: error while trying to generate symmetric key", "error", err)
	}
	if err := app.LoadEnv(".env"); err != nil {
		logger.Error("main:error while trying to read env", "error", err)
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://admin:root@localhost:5432/bazaargo?sslmode=disable"
	}
	logger.Info("trying to connect to db")
	log.Println(dbURL)
	db, err := postgres.NewConn(dbURL, 20, 10, time.Minute)
	if err != nil {
		logger.Error("main: error while connecting to db", "error", err)
		os.Exit(1)
	}
	logger.Info("db connection successful")
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("main: error while closing db conn", "error", err)
		}
	}()
	logger.Info("trying to start server")

	app, err := app.New(db, logger, symmetricKey)
	if err != nil {
		logger.Error("main:error while building app", "error", err)
		os.Exit(1)
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	if err := app.Handler.Run(app.Handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")

}
