package logger

import (
	"go.uber.org/zap"
	"log"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}
	defer logger.Sync()

	return &Logger{logger.Sugar()}
}