package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates a new logger instance
func New(level string) (*logrus.Logger, error) {
	log := logrus.New()

	// Parse log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(lvl)

	// Set output
	log.SetOutput(os.Stdout)

	// Set format
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Enable caller info
	log.SetReportCaller(true)

	return log, nil
} 