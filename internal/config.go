package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"application-design/pkg/environment"
)

type Config struct {
	Server Server
	Sentry Sentry
}

func LoadConfig() (Config, error) {
	var (
		config Config
		err    error
	)

	config.Server, err = loadServer()
	if err != nil {
		return config, fmt.Errorf("could not load server config: %w", err)
	}
	config.Sentry = loadSentry()

	return config, nil
}

type Server struct {
	Environment      environment.Type
	Port             string
	InterruptTimeout time.Duration
	PprofPort        string
}

func loadServer() (Server, error) {
	var server Server

	server.Environment = environment.Type(getEnv("ENV_TYPE", string(environment.Testing)))
	server.Port = getEnv("SERVER_PORT", "8080")
	interruptTimeout, err := time.ParseDuration(getEnv("KILL_TIMEOUT", "2s"))
	if err != nil {
		return server, fmt.Errorf("could not parse kill timeout: %w", err)
	}
	server.InterruptTimeout = interruptTimeout
	server.PprofPort = getEnv("PPROF_PORT", "6060")

	return server, nil
}

type Sentry struct {
	DSN         string
	Environment environment.Type
}

func loadSentry() Sentry {
	var sentry Sentry

	sentry.Environment = environment.Type(getEnv("SENTRY_ENVIRONMENT", string(environment.Testing)))
	sentry.DSN = getEnv("SENTRY_DSN", "")

	return sentry
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}
	return fallback
}
