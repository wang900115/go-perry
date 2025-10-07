package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(file), zap.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	core := zapcore.NewTee(fileCore, consoleCore)
	logger := zap.New(core)

	defer logger.Sync()

	logger.Info("測試日誌", zap.String("user", "Alice"))
}

// func main() {
// 	config := zap.Config{
// 		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
// 		Development: false,
// 		Encoding:    "json",
// 		EncoderConfig: zapcore.EncoderConfig{
// 			TimeKey:        "time",
// 			LevelKey:       "level",
// 			NameKey:        "logger",
// 			CallerKey:      "caller",
// 			MessageKey:     "msg",
// 			StacktraceKey:  "stacktrace",
// 			LineEnding:     zapcore.DefaultLineEnding,
// 			EncodeLevel:    zapcore.LowercaseLevelEncoder,
// 			EncodeTime:     zapcore.ISO8601TimeEncoder,
// 			EncodeDuration: zapcore.SecondsDurationEncoder,
// 			EncodeCaller:   zapcore.ShortCallerEncoder,
// 		},
// 		OutputPaths:      []string{"stdout", "log.txt"},
// 		ErrorOutputPaths: []string{"stderr"},
// 	}

// 	logger, err := config.Build()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer logger.Sync()

// 	logger.Info("自定義配置的文件", zap.String("key", "value"))
// }
