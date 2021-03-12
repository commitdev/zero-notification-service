package config

import (
	"strings"
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
	AllowEmailToDomains     []string
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
	AllowEmailToDomains
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

	viper.SetDefault(AllowEmailToDomains, []string{})
	viper.BindEnv(AllowEmailToDomains, "ALLOW_EMAIL_TO_DOMAINS")

	// Split the string on commas, viper doesn't support doing this to env vars
	domains := []string{}
	if strings.Trim(viper.GetString(AllowEmailToDomains), " ") != "" {
		domains = strings.Split(viper.GetString(AllowEmailToDomains), ",")
		for i, domain := range domains {
			domains[i] = strings.Trim(domain, " ")
		}
	}

	config := Config{
		Port:                    viper.GetInt(Port),
		SendgridAPIKey:          viper.GetString(SendgridAPIKey),
		SlackAPIKey:             viper.GetString(SlackAPIKey),
		GracefulShutdownTimeout: viper.GetDuration(GracefulShutdownTimeout),
		StructuredLogging:       viper.GetBool(StructuredLogging),
		DebugDumpRequests:       viper.GetBool(DebugDumpRequests),
		AllowEmailToDomains:     domains,
	}

	return &config
}
