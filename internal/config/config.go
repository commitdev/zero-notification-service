package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config is the struct that holds configuration values
type Config struct {
	Port                    int
	SendgridAPIKey          string
	SlackAPIKey             string
	GracefulShutdownTimeout time.Duration
	StructuredLogging       bool
	DebugDumpRequests       bool
}

var config *Config

const (
	_ = string(rune(iota)) // We don't care about the values of these constants
	Port
	SendgridAPIKey
	SlackAPIKey
	GracefulShutdownTimeout
	StructuredLogging
	DebugDumpRequests
)

// GetConfig returns a pointer to the singleton Config object
func GetConfig() *Config {
	if config == nil {
		config = loadConfig()
	}
	return config
}

func loadConfig() *Config {
	viper.SetDefault(Port, "80")
	viper.BindEnv(Port, "SERVICE_PORT")

	viper.SetDefault(SendgridAPIKey, "")
	viper.BindEnv(SendgridAPIKey, "SENDGRID_API_KEY")

	viper.SetDefault(SlackAPIKey, "")
	viper.BindEnv(SlackAPIKey, "SLACK_API_KEY")

	viper.SetDefault(GracefulShutdownTimeout, "10")
	viper.BindEnv(GracefulShutdownTimeout, "GRACEFUL_SHUTDOWN_TIMEOUT_SECONDS")

	viper.SetDefault(StructuredLogging, "false")
	viper.BindEnv(StructuredLogging, "STRUCTURED_LOGGING")

	viper.SetDefault(DebugDumpRequests, "false")
	viper.BindEnv(DebugDumpRequests, "DEBUG_DUMP_REQUESTS")

	config := Config{
		Port:                    viper.GetInt(Port),
		SendgridAPIKey:          viper.GetString(SendgridAPIKey),
		SlackAPIKey:             viper.GetString(SlackAPIKey),
		GracefulShutdownTimeout: viper.GetDuration(GracefulShutdownTimeout),
		StructuredLogging:       viper.GetBool(StructuredLogging),
		DebugDumpRequests:       viper.GetBool(DebugDumpRequests),
	}

	return &config
}
