package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/flagforge/flagforge/worker/internal/config"
	"github.com/flagforge/flagforge/worker/internal/jobs"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.With().Str("service", "worker").Logger()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	consumer := jobs.NewAuditConsumer(rdb, logger)
	go consumer.Run(ctx)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		logger.Info().Int("port", cfg.Port).Msg("worker listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("worker server error")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	cancel()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTTL)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("worker shutdown failed")
	}

	logger.Info().Msg("worker shutdown complete")
}
