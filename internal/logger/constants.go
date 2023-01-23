package logger

import "go.uber.org/zap/zapcore"

const Encoding = "json"

var EncoderConfig = zapcore.EncoderConfig{
	MessageKey:  "message",
	LevelKey:    "level",
	EncodeLevel: zapcore.CapitalLevelEncoder,
	TimeKey:     "time",
	EncodeTime:  zapcore.ISO8601TimeEncoder,
}
