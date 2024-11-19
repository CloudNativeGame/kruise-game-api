package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var opts = &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))

func init() {
	slog.SetDefault(logger)
}

func main() {
	r := gin.New()

	config := sloggin.Config{
		WithRequestBody:    true,
		WithResponseBody:   true,
		WithRequestHeader:  true,
		WithResponseHeader: true,
		Filters: []sloggin.Filter{
			sloggin.IgnorePath("/healthz"),
		},
	}
	r.Use(sloggin.NewWithConfig(logger, config))
	r.Use(gin.Recovery())
	registerRoutes(r)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve failed", "error", err)
		}
	}()
	gracefulShutdown(&srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shut down", "error", err)
		panic(err)
	}
}
