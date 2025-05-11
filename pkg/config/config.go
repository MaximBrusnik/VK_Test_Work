package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig   `json:"log"`
	PubSub PubSubConfig `json:"pubsub"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	// Port is the port number the server will listen on
	Port int `json:"port" validate:"required,min=1,max=65535"`
	
	// Host is the host address the server will listen on
	Host string `json:"host" validate:"required"`
	
	// GracefulShutdownTimeout is the maximum time to wait for graceful shutdown
	GracefulShutdownTimeout time.Duration `json:"graceful_shutdown_timeout" validate:"required"`
	
	// MaxConcurrentStreams is the maximum number of concurrent streams per connection
	MaxConcurrentStreams uint32 `json:"max_concurrent_streams" validate:"required,min=1"`
}

// LogConfig contains logging-related configuration
type LogConfig struct {
	// Level is the logging level (debug, info, warn, error, fatal, panic)
	Level string `json:"level" validate:"required,oneof=debug info warn error fatal panic"`
	
	// Format is the log format (json or text)
	Format string `json:"format" validate:"required,oneof=json text"`
	
	// Output is the log output (stdout, stderr, or file path)
	Output string `json:"output" validate:"required"`
	
	// EnableCaller enables logging of caller information
	EnableCaller bool `json:"enable_caller"`
	
	// TimestampFormat is the format for timestamps in logs
	TimestampFormat string `json:"timestamp_format"`
}

// PubSubConfig contains pubsub-related configuration
type PubSubConfig struct {
	// MaxSubscribersPerKey is the maximum number of subscribers per key
	MaxSubscribersPerKey int `json:"max_subscribers_per_key" validate:"required,min=1"`
	
	// MessageBufferSize is the size of the message buffer for each subscriber
	MessageBufferSize int `json:"message_buffer_size" validate:"required,min=1"`
	
	// CleanupInterval is the interval for cleaning up inactive subscribers
	CleanupInterval time.Duration `json:"cleanup_interval" validate:"required"`
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                  8080,
			Host:                  "0.0.0.0",
			GracefulShutdownTimeout: 30 * time.Second,
			MaxConcurrentStreams:  100,
		},
		Log: LogConfig{
			Level:           "info",
			Format:          "json",
			Output:          "stdout",
			EnableCaller:    true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
		PubSub: PubSubConfig{
			MaxSubscribersPerKey: 1000,
			MessageBufferSize:    100,
			CleanupInterval:      5 * time.Minute,
		},
	}
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}

	if c.Server.GracefulShutdownTimeout <= 0 {
		return fmt.Errorf("graceful shutdown timeout must be positive")
	}

	if c.Server.MaxConcurrentStreams < 1 {
		return fmt.Errorf("max concurrent streams must be positive")
	}

	switch c.Log.Level {
	case "debug", "info", "warn", "error", "fatal", "panic":
	default:
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	switch c.Log.Format {
	case "json", "text":
	default:
		return fmt.Errorf("invalid log format: %s", c.Log.Format)
	}

	if c.Log.Output == "" {
		return fmt.Errorf("log output is required")
	}

	if c.PubSub.MaxSubscribersPerKey < 1 {
		return fmt.Errorf("max subscribers per key must be positive")
	}

	if c.PubSub.MessageBufferSize < 1 {
		return fmt.Errorf("message buffer size must be positive")
	}

	if c.PubSub.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup interval must be positive")
	}

	return nil
} 