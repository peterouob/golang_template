package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logg *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)
	//
	//logFile := getLogFile()
	//writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile))
	//
	//core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logg = zap.New(core, zap.AddCaller())
}

//func getLogFile() *os.File {
//	today := time.Now().Format("2006-01-02")
//	logFileName := fmt.Sprintf("./tools/log/log-%s.log", today)
//
//	err := os.MkdirAll("./tools/log/", os.ModePerm)
//	if err != nil {
//		panic(fmt.Sprintf("failed to create log directory: %v", err))
//	}
//
//	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
//	if err != nil {
//		panic(fmt.Sprintf("failed to open log file: %v", err))
//	}
//	return file
//}

func getLogger() *zap.Logger {
	if logg == nil {
		InitLogger()
	}
	return logg
}

func Log(msg any) {
	getLogger().Info(msg.(string))
}

func Logf(format string, args ...any) {
	getLogger().Info(fmt.Sprintf(format, args...))
}

func Warn(msg any) {
	getLogger().Warn(msg.(string))
}

func Error(msg string, err error, fields ...zap.Field) {
	getLogger().Error(msg, append(fields, zap.Error(err))...)
}

func HandelError(msg string, err error, f ...func(args ...any)) {
	if err != nil {
		if len(f) > 0 && f[0] != nil {
			f[0](msg, err)
		}

		Error(msg, err)
	}
}

func ErrorMsg(msg string) {
	getLogger().Error(msg)
}

func ErrorMsgF(format string, args ...any) {
	getLogger().Error(fmt.Sprintf(format, args...))
}

func FormatString(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
