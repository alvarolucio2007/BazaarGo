package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/app"
	httpadapter "github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/inbound/http"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/outbound/gemini"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/adapters/outbound/pgvector"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/service"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	if err := app.LoadEnv(".env"); err != nil {
		log.Fatalf("main:error while trying to read env: %v", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbURL := os.Getenv("DATABASE_LOCAL_URL")
	if dbURL == "" {
		dbURL = "postgresql://admin:root@localhost:5432/bazaargo?sslmode=disable"
	}

	db, err := pgvector.NewConn(dbURL, 20, 10, time.Minute)
	if err != nil {
		log.Fatalf("error while connecting to DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error while closing DB: %v", err)
		}
	}()

	// 3. Instanciando o VectorStore (Faltava essa linha!)
	pgStore := pgvector.NewVectorStore(db, logger)

	// 4. Configurando o Gemini
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("api key empty bro check yo .env")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("error while creating genai client: %v", err)
	}

	// 5. Instanciando o Adapter do Gemini
	geminiAdapter := gemini.NewGeminiAdapter(client)

	// 6. Injetando tudo no Service
	chatSvc := service.NewChatService(geminiAdapter, pgStore, geminiAdapter)

	// 7. Configurando o Handler HTTP
	handler := httpadapter.NewChatHandler(chatSvc) // Assumindo que você criou esse construtor no pacote http

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8081"
	}
	log.Println("starting server yo")
	if err := handler.Run(handler.Mount(), addr); err != nil {
		logger.Error("main: error while booting HTTP server", "error", err)
		os.Exit(1)
	}
	logger.Info("server has been shut down")
}
