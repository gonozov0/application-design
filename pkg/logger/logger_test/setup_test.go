package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"application-design/pkg/logger"
)

func TestSetup(_ *testing.T) {
	logger.Setup()
	slog.InfoContext(context.Background(), "test logging")
}
