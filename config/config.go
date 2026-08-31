package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

type Config struct {
	Environment Environment
	Host        string
	Port        string
	LogLevel    slog.Level

	// TrustProxy says whether X-Forwarded-* headers on an inbound request can
	// be believed. Only turn it on when something in front of this server
	// (Cloudflare, Fly, a reverse proxy) actually sets them: a client can put
	// any value it likes in those headers, so trusting them on a directly
	// exposed server hands out both IP spoofing and rate-limit evasion.
	TrustProxy bool

	// ClientIPHeader names the single header the proxy in front is known to
	// *overwrite* -- CF-Connecting-IP behind a Cloudflare tunnel. It is read
	// only when TrustProxy is on, so naming a header can never by itself make
	// a directly exposed server believe one.
	//
	// The header has to be one the proxy replaces rather than appends to.
	// X-Forwarded-For is the trap: Cloudflare appends the real address to
	// whatever the client already put there, so its leftmost value is the
	// client's own and trusting it is the same as trusting nothing.
	ClientIPHeader string

	// HostChassis is the machine named in the footer beside the CPU the
	// process actually reads from /proc. Nothing inside an unprivileged
	// container can identify the box, so it has to be stated -- but stating it
	// in the binary means it follows that binary onto every other machine it
	// ever runs on, and a hardcoded chassis beside a live CPU model is exactly
	// the pairing nobody checks. Setting it per deployment keeps it true where
	// it is set and absent where it is not; empty drops the segment.
	HostChassis string
}

var (
	Global *Config
	once   sync.Once
)

func init() {
	once.Do(func() {
		Global = Load()
	})
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadBase() *Config {
	godotenv.Load()

	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),
		LogLevel: func() slog.Level {
			switch os.Getenv("LOG_LEVEL") {
			case "DEBUG":
				return slog.LevelDebug
			case "INFO":
				return slog.LevelInfo
			case "WARN":
				return slog.LevelWarn
			case "ERROR":
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(),
		TrustProxy:     getEnv("TRUST_PROXY", "") == "true",
		ClientIPHeader: getEnv("CLIENT_IP_HEADER", "CF-Connecting-IP"),
		HostChassis:    getEnv("HOST_CHASSIS", ""),
	}
}
