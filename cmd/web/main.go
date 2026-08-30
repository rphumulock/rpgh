package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rpgh/config"
	"rpgh/router"
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

	// Where a request's client address comes from, decided once and explicitly.
	// The TCP peer is the only source nothing can forge, so it is the default;
	// a header is believed only when something in front is known to overwrite
	// it. This has to be installed before the rate limiter, which keys off it.
	if config.Global.TrustProxy && config.Global.ClientIPHeader != "" {
		r.Use(middleware.ClientIPFromHeader(config.Global.ClientIPHeader))
	} else {
		r.Use(middleware.ClientIPFromRemoteAddr)
	}

	r.Use(
		httplog.RequestLogger(logger, nil),
		middleware.Recoverer,
		router.SecurityHeaders,
		// A page load costs about ten requests, so this leaves room for normal
		// browsing while capping what a single address can pull. It protects
		// CPU against one noisy source; it is not a defence against a
		// distributed flood, which has to be absorbed upstream.
		httprate.LimitBy(240, time.Minute, clientIPKey),
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

// clientIPKey buckets the rate limiter by the address resolved above rather
// than by a header read at limit time. CanonicalizeIP groups IPv6 by /64, so a
// single client cannot walk its own prefix for a fresh bucket per request.
func clientIPKey(r *http.Request) (string, error) {
	if ip := middleware.GetClientIP(r.Context()); ip != "" {
		return httprate.CanonicalizeIP(ip), nil
	}

	// ClientIPFromHeader fails closed: a request that did not carry the
	// trusted header leaves no address behind. That is every caller which
	// reached us without going through the edge -- a reverse proxy on the LAN,
	// a health probe -- and keying them all on the empty string would put them
	// in one shared bucket, where they would rate-limit each other. The peer
	// is the honest answer for those, and is not forgeable.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return httprate.CanonicalizeIP(host), nil
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
