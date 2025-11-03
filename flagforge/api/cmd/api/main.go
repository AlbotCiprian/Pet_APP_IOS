package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/flagforge/flagforge/api/internal/config"
	"github.com/flagforge/flagforge/api/internal/handlers"
	httptransport "github.com/flagforge/flagforge/api/internal/http"
	"github.com/flagforge/flagforge/api/internal/repo"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.With().Str("service", "api").Logger()

	ctx := context.Background()

	db := connectPostgres(ctx, cfg.PostgresDSN, logger)
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	if err := runMigrations(ctx, db, cfg.Migrations); err != nil {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}

	store := repo.NewStore(db, rdb, logger)
	h := handlers.Handler{Store: store, Log: logger}

	r := httptransport.NewRouter(h)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		logger.Info().Int("port", cfg.Port).Msg("api listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server error")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTTL)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}

	logger.Info().Msg("shutdown complete")
}

func connectPostgres(ctx context.Context, dsn string, logger zerolog.Logger) *pgxpool.Pool {
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			timeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			e := pool.Ping(timeCtx)
			cancel()
			if e == nil {
				return pool
			}
			err = e
		}
		logger.Warn().Err(err).Int("attempt", i+1).Msg("waiting for postgres")
		time.Sleep(2 * time.Second)
	}
	logger.Fatal().Err(err).Msg("unable to connect to postgres")
	return nil
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := conn.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %w", entry.Name(), err)
		}
	}

	return nil
}
