package initialize

import (
	"fmt"
	"iptv-spider-sh/global"
	"iptv-spider-sh/utils"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AccessLog() *zap.Logger {
	cfg := global.CONFIG.AccessLog
	if !cfg.Enabled {
		return nil
	}

	if ok, _ := utils.PathExists(cfg.Director); !ok {
		fmt.Printf("create %v directory\n", cfg.Director)
		_ = os.MkdirAll(cfg.Director, os.ModePerm)
	}

	encoder := getAccessLogEncoder(cfg.Format)
	writer, err := utils.GetAccessLogWriteSyncer(cfg)
	if err != nil {
		fmt.Printf("Get Access Log Write Syncer Failed err:%v", err.Error())
		writer = zapcore.AddSync(os.Stdout)
	}

	if cfg.LogInConsole {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer)
	}

	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	return zap.New(core)
}

func getAccessLogEncoder(format string) zapcore.Encoder {
	config := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     accessLogTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	if format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

func accessLogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[ACCESS] 2006-01-02 15:04:05.000"))
}
