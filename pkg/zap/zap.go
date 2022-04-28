package zap

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(level, format, filepath string) (*Logger, error) {
	var logLevel zapcore.Level

	defer func() {
		if err := recover(); err != nil {
			log.Println("zap: zap: panic occurred during zap initialization:", err)
		}
	}()

	switch strings.ToLower(level) {
	case "error":
		logLevel = zapcore.ErrorLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "debug":
		logLevel = zapcore.DebugLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	var fileEncoder zapcore.Encoder

	switch strings.ToLower(format) {
	case "json":
		fileEncoder = getJSONEncoder()
	case "console":
		fileEncoder = getConsoleEncoder()
	default:
		fileEncoder = getJSONEncoder()
	}

	consoleEncoder := getConsoleEncoder()

	writeSyncer, err := getLogWriter(filepath)
	if err != nil {
		return nil, fmt.Errorf("zap.New: ")
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writeSyncer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel))

	logger := zap.New(core, zap.AddCaller())

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println(err)
		}
	}()

	return &Logger{logger.Sugar()}, nil
}

func getConsoleEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(config)
}

func getJSONEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(config)
}

func getLogWriter(filePath string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0744)
	if err != nil {
		return nil, fmt.Errorf("getLogWriter: failed to open file: %w", err)
	}

	return zapcore.AddSync(file), nil
}
