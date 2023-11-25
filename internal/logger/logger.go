package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chirikova/go-anti-brute-force/internal/config"
	log "github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
}

type AppLogger struct {
	logger *log.Logger
}

func New(cfg config.Logger, output io.Writer) (Logger, error) {
	logLevel, err := log.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	logger := log.New()
	logger.SetLevel(logLevel)
	logger.SetFormatter(
		&log.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "02/Jan/2006 15:04:05 -07:00",
			DisableLevelTruncation: true,
		},
	)
	logger.SetOutput(output)

	return &AppLogger{logger}, nil
}

func GetLogFile(filePath string) (*os.File, error) {
	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0o700)
		if err != nil {
			return nil, fmt.Errorf("create dir %s: %w", dirPath, err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *AppLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *AppLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
