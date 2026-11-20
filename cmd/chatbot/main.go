package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/logger"
	"github.com/alvarolucio2007/bazaargo/internal/app"
	httpadapter "github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/inbound/http"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/outbound/gemini"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/outbound/pgvector"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/service"
	"github.com/alvarolucio2007/bazaargo/internal/telemetry"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	if err := app.LoadEnv(".env"); err != nil {
		log.Fatalf("main:error while trying to read env: %v", err)
	}
	otelLoggerAddr := os.Getenv("OTEL_LOGGER_ENDPOINT")
	if otelLoggerAddr == "" {
		otelLoggerAddr = "localhost:4318"
	}
	lp, err := telemetry.InitLogger(context.Background(), otelLoggerAddr)
	if err != nil {
		panic("failed to initialize logger provider: " + err.Error())
	}
	defer func() {
		if err := lp.Shutdown(context.Background()); err != nil {
			log.Println("failed to shutdown logger provider:", err)
		}
	}()
	logger := logger.NewSlog()
	dbURL := os.Getenv("DATABASE_LOCAL_URL")
	if dbURL == "" {
		dbURL = "postgresql://admin:root@localhost:5432/bazaargo?sslmode=disable"
	}

	db, err := pgvector.NewConn(dbURL, 20, 10, time.Minute)
	if err != nil {
		logger.Error("error while connecting to DB", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error while closing DB", "error", err)
			os.Exit(1)
		}
	}()

	pgStore := pgvector.NewVectorStore(db, logger)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Error("Gemini API key is empty")
		os.Exit(1)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		logger.Error("error while creating genai client", "error", err)
		os.Exit(1)
	}

	geminiAdapter := gemini.NewGeminiAdapter(client)

	chatSvc := service.NewChatService(geminiAdapter, pgStore, geminiAdapter)

	handler := httpadapter.NewChatHandler(chatSvc) // Assumindo que você criou esse construtor no pacote http

	logger.Info("starting telemetry")
	otelAddr := os.Getenv("OTEL_ENDPOINT")
	if otelAddr == "" {
		otelAddr = "localhost:4317"
	}
	tp, err := telemetry.InitTracer(context.Background(), otelAddr)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("failed to shutdown tracer", "error", err)
		}
	}()
	mp, err := telemetry.InitMetrics(context.Background(), otelAddr)
	if err != nil {
		logger.Error("failed to initialize metrics", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mp.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown metrics", "error", err)
		}
	}()

	logger.Info("tracer started successfully")
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	logger.Info("starting server")
	if err := handler.Run(handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")
}
