package tools

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var logg *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)

	logFile := getLogFile()
	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile))

	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	logg = zap.New(core, zap.AddCaller())
}

func getLogFile() *os.File {
	today := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("./tools/log/log-%s.log", today)

	err := os.MkdirAll("./tools/log/", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
	return file
}

func LogMessage(msg string) {
	if logg == nil {
		panic("logger is not initialized")
	}
	logg.Info(msg)
}

func LogError(msg string, err error) {
	if logg == nil {
		panic("logger is not initialized")
	}
	logg.Error(msg, zap.Error(err))
}
