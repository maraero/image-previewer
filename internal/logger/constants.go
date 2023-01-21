package logger

import "go.uber.org/zap/zapcore"

const (
	Level    = "info"
	Encoding = "json"
)

var (
	OutputPaths      = []string{"./logs/log.log", "stdout"}
	ErrorOutputPaths = []string{"./logs/log.log", "stderr"}
)

var EncoderConfig = zapcore.EncoderConfig{
	MessageKey:  "message",
	LevelKey:    "level",
	EncodeLevel: zapcore.CapitalLevelEncoder,
	TimeKey:     "time",
	EncodeTime:  zapcore.ISO8601TimeEncoder,
}
