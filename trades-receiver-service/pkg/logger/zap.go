package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println(err)
		}
	}()

	return &Logger{logger.Sugar()}
}
