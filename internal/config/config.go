package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config is the struct that holds configuration values
type Config struct {
	Port                    int
	SendgridAPIKey          string
	GracefulShutdownTimeout time.Duration
}

var config *Config

const (
	_ = string(iota) // We don't care about the values of these constants
	Port
	SendgridAPIKey
	GracefulShutdownTimeout
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

	viper.SetDefault(GracefulShutdownTimeout, "10")
	viper.BindEnv(GracefulShutdownTimeout, "GRACEFUL_SHUTDOWN_TIMEOUT_SECONDS")

	config := Config{
		Port:                    viper.GetInt(Port),
		SendgridAPIKey:          viper.GetString(SendgridAPIKey),
		GracefulShutdownTimeout: viper.GetDuration(GracefulShutdownTimeout),
	}

	return &config
}
