package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task/internal/config"
	flood_control "task/internal/flood-control"
	"task/internal/usecase"
	"time"
)

func Run() {
	cfg := config.ConfgiLoad()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	usersRate := usecase.NewUserRateUseCase(client)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Info("failed to connect to redis", slog.Any("error", error.Error))
		return
	}

	defer client.Close()

	router := http.NewServeMux()

	handler := flood_control.NewHandler(usersRate, cfg.RequestCount, cfg.Interval)
	router.Handle("/", handler.Nothing(context.Background(), log))

	log.Info("successful connect to redis")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		srv.ListenAndServe()
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Any("error", err.Error))
		return
	}

	log.Info("server stopped")
}
