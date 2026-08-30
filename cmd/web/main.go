package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"rpgh/config"
	"rpgh/router"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("server failure", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.Global.LogLevel,
	}))
	slog.SetDefault(logger)

	r := chi.NewMux()

	// RealIP rewrites RemoteAddr from X-Forwarded-For, so it must only run when
	// that header comes from a proxy we trust -- otherwise every client picks
	// its own identity and the rate limiter below counts a spoofed key.
	if config.Global.TrustProxy {
		r.Use(middleware.RealIP)
	}

	r.Use(
		httplog.RequestLogger(logger, nil),
		middleware.Recoverer,
		router.SecurityHeaders,
		// A page load costs about ten requests, so this leaves room for normal
		// browsing while capping what a single address can pull. It protects
		// CPU against one noisy source; it is not a defence against a
		// distributed flood, which has to be absorbed upstream.
		httprate.LimitByIP(240, time.Minute),
		middleware.Compress(5),
	)

	eg, egctx := errgroup.WithContext(ctx)

	if err := router.SetupRoutes(r); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,

		// Without these a connection can hold a goroutine open indefinitely:
		// Slowloris is just many connections dribbling a header a byte at a
		// time, which costs the client nothing and needs no volume.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      writeTimeout(),
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 16,

		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		ErrorLog: slog.NewLogLogger(
			slog.Default().Handler(),
			slog.LevelError,
		),
	}

	eg.Go(func() error {
		slog.Info("server started", "addr", srv.Addr)
		err := srv.ListenAndServe()

		if err == nil || err == http.ErrServerClosed {
			return nil
		}

		return fmt.Errorf("server error: %w", err)
	})

	eg.Go(func() error {
		<-egctx.Done()

		slog.Debug("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Debug("shutting down server...")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		return nil
	})

	return eg.Wait()
}

// writeTimeout bounds how long a response may take to write. It has to stay
// unset in dev: the hot-reload endpoint holds an SSE stream open until the
// next rebuild, and a write deadline would cut it every time it expired.
func writeTimeout() time.Duration {
	if config.Global.Environment == config.Dev {
		return 0
	}
	return 30 * time.Second
}
