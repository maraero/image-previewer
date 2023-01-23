package logger

import (
	"fmt"

	"github.com/maraero/image-previewer/internal/config"
	"go.uber.org/zap"
)

func buildZapConfig(cfg config.Logger) (zap.Config, error) {
	parsedLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return zap.Config{}, fmt.Errorf("can not parse level: %w", err)
	}
	return zap.Config{
		Encoding:         Encoding,
		Level:            parsedLevel,
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
		EncoderConfig:    EncoderConfig,
	}, nil
}
