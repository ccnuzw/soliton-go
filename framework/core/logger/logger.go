package logger

import (
	"os"
	"strings"

	"github.com/soliton-go/framework/core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap Logger.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// Simple production-style config by default.
	// We still honor config-driven log level for quick tuning.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zap.InfoLevel
	switch strings.ToLower(strings.TrimSpace(cfg.GetString("log.level"))) {
	case "debug":
		level = zap.DebugLevel
	case "warn", "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)
	
	return zap.New(core), nil
}
