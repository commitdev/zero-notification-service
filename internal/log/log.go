package log

import (
	"log"

	"github.com/commitdev/zero-notification-service/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init sets up logging based on the current environment
func Init(config *config.Config) {
	var rawLogger *zap.Logger
	var err error
	if config.StructuredLogging {
		// Info level, JSON output
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.MessageKey = "message"
		zapConfig.EncoderConfig.LevelKey = "log.level"
		zapConfig.EncoderConfig.CallerKey = "log.origin.file.name"
		zapConfig.EncoderConfig.TimeKey = "@timestamp"
		zapConfig.EncoderConfig.FunctionKey = "log.origin.function"
		rawLogger, err = zapConfig.Build()
	} else {
		// Debug level, pretty output
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		rawLogger, err = zapConfig.Build()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	zap.ReplaceGlobals(rawLogger)
}
