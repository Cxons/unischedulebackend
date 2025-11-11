package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authHandler "github.com/Cxons/unischedulebackend/internal/auth/handler"
	regHandler "github.com/Cxons/unischedulebackend/internal/registration/handler"
	"github.com/Cxons/unischedulebackend/internal/server"
	"github.com/Cxons/unischedulebackend/internal/shared/db"
	"github.com/Cxons/unischedulebackend/pkg/caching"
	supHandler "github.com/Cxons/unischedulebackend/pkg/supabase/handler"
	"github.com/joho/godotenv"
)

func main() {
	// 1️⃣ Load environment variables from .env file
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("⚠️ No .env file found, relying on system environment variables")
	}

	// 2️⃣ Setup JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger.Info("🚀 Starting UniSchedule backend...")

	// 3️⃣ Load configuration from environment
	cfg := server.Config{
		PORT: getEnv("PORT", "5000"),
		ENV:  getEnv("ENV", "development"),
	}

	// 4️⃣ Initialize PostgreSQL database
	dbInstance, err := db.NewDatabase(getEnv("DATABASE_URL", ""))
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("✅ Connected to PostgreSQL")

	// 5️⃣ Initialize Redis client
	cacheClient := caching.NewRedisClient(
		getEnv("REDIS_ADDR", "localhost:6379"),
		getEnv("REDIS_PASSWORD", ""),
		0, // DB index
		2, // Protocol (2 = RESP2, 3 = RESP3)
	)
	logger.Info("✅ Connected to Redis")

	// 6️⃣ Initialize module handlers
	auth := authHandler.NewAuthPackage(logger, dbInstance.DB)
	reg := regHandler.NewRegPackage(logger, dbInstance.DB)
	supabase := supHandler.NewSupabasePackage(logger,getEnv("SUPABASE_URL",""),getEnv("SUPABASE_SECRET_KEY",""))


	// 7️⃣ Create HTTP server
	srv := server.NewServer(cfg, logger, cacheClient, dbInstance, auth, reg,supabase)

	// 8️⃣ Start the server in a goroutine
	go func() {
		logger.Info("🌐 Server running", slog.String("address", ":"+cfg.PORT))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 9️⃣ Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Shutdown signal received, cleaning up resources...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop HTTP server gracefully
	if err := srv.Server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown failed", slog.String("error", err.Error()))
	} else {
		logger.Info("✅ HTTP server stopped gracefully")
	}

	// Close DB connection
	dbInstance.CloseConnection()
	logger.Info("✅ PostgreSQL connection closed")

	// Close Redis connection
	cacheClient.Db.Close()
	logger.Info("✅ Redis connection closed")

	logger.Info("✨ UniSchedule backend stopped successfully")
}

// getEnv returns the value of an environment variable or a default fallback.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
