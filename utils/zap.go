package utils

import (
	"iptv-spider-sh/config"
	"iptv-spider-sh/global"
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(global.CONFIG.Zap.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(global.CONFIG.Zap.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if global.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

func GetAccessLogWriteSyncer(cfg config.AccessLog) (zapcore.WriteSyncer, error) {
	maxAge := 7 * 24 * time.Hour
	if cfg.MaxAge > 0 {
		maxAge = time.Duration(cfg.MaxAge) * 24 * time.Hour
	}
	rotationTime := 24 * time.Hour
	if cfg.RotationTime > 0 {
		rotationTime = time.Duration(cfg.RotationTime) * time.Hour
	}

	fileWriter, err := zaprotatelogs.New(
		path.Join(cfg.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(maxAge),
		zaprotatelogs.WithRotationTime(rotationTime),
	)
	return zapcore.AddSync(fileWriter), err
}
