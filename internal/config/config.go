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
	TwilioAccID             string
	TwilioAuthToken         string
	TwilioPhoneNumber       string
	GracefulShutdownTimeout time.Duration
	StructuredLogging       bool
}

var config *Config

const (
	_ = string(rune(iota)) // We don't care about the values of these constants
	Port
	SendgridAPIKey
	SlackAPIKey
	TwilioAccID
	TwilioAuthToken
	TwilioPhoneNumber
	GracefulShutdownTimeout
	StructuredLogging
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

	viper.SetDefault(TwilioAccID, "")
	viper.BindEnv(TwilioAccID, "TWILIO_ACC_ID")

	viper.SetDefault(TwilioAuthToken, "")
	viper.BindEnv(TwilioAuthToken, "TWILIO_AUTH_TOKEN")

	viper.SetDefault(TwilioPhoneNumber, "")
	viper.BindEnv(TwilioPhoneNumber, "TWILIO_PHONE_NUMBER")

	viper.SetDefault(GracefulShutdownTimeout, "10")
	viper.BindEnv(GracefulShutdownTimeout, "GRACEFUL_SHUTDOWN_TIMEOUT_SECONDS")

	viper.SetDefault(StructuredLogging, "false")
	viper.BindEnv(StructuredLogging, "STRUCTURED_LOGGING")

	config := Config{
		Port:                    viper.GetInt(Port),
		SendgridAPIKey:          viper.GetString(SendgridAPIKey),
		SlackAPIKey:             viper.GetString(SlackAPIKey),
		TwilioAccID:             viper.GetString(TwilioAccID),
		TwilioAuthToken:         viper.GetString(TwilioAuthToken),
		TwilioPhoneNumber:       viper.GetString(TwilioPhoneNumber),
		GracefulShutdownTimeout: viper.GetDuration(GracefulShutdownTimeout),
		StructuredLogging:       viper.GetBool(StructuredLogging),
	}

	return &config
}
