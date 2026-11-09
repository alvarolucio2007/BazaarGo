package main

import (
	"context"
	"os"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/logger"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/mongo"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/postgres"
	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/valkey"
	"github.com/alvarolucio2007/bazaargo/internal/app"
	"github.com/alvarolucio2007/bazaargo/internal/telemetry"
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

	logger.Info("trying to connect to mongodb")

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:admin123@localhost:27017/bazaargo?authSource=admin"
	}
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		mongoDBName = "bazaargo"
	}
	mongoCollname := os.Getenv("MONGO_COLL_NAME")
	if mongoCollname == "" {
		mongoCollname = "products"
	}
	mdb, err := mongo.NewConn(mongoURL)
	if err != nil {
		logger.Error("main: error while connecting to mongo", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := mdb.Disconnect(context.Background()); err != nil {
			logger.Error("main: error while closing mongo conn", "error", err)
		}
	}()

	mColl := mongo.ConnToCollection(mdb, mongoDBName, mongoCollname)

	logger.Info("mongo connection successful")

	logger.Info("trying to connect to cache")

	valkeyURL := os.Getenv("CACHE_URL")
	if valkeyURL == "" {
		valkeyURL = "localhost:6379"
	}
	rdb, err := valkey.NewConn(valkeyURL, "", 0)
	if err != nil {
		logger.Error("main: error while connecting to cache", "error", err)
		os.Exit(1)
	}

	app, err := app.New(db, mColl, rdb, logger, symmetricKey)
	if err != nil {
		logger.Error("main:error while building app", "error", err)
		os.Exit(1)
	}
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	logger.Info("starting telemetry")
	tp, err := telemetry.InitTracer(context.Background(), "bazaargo-api")
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("failed to shut down tracer", "error", err)
			os.Exit(1)
		}
	}()
	logger.Info("telemetry started sucessfully")
	logger.Info("starting server")

	if err := app.Handler.Run(app.Handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")

}
