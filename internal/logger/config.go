package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func buildZapConfig() (zap.Config, error) {
	parsedLevel, err := zap.ParseAtomicLevel(Level)
	if err != nil {
		return zap.Config{}, fmt.Errorf("can not parse level: %w", err)
	}
	return zap.Config{
		Encoding:         Encoding,
		Level:            parsedLevel,
		OutputPaths:      OutputPaths,
		ErrorOutputPaths: ErrorOutputPaths,
		EncoderConfig:    EncoderConfig,
	}, nil
}
