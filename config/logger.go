package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(debug bool) (*zap.Logger, error) {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	jsonEnc := zapcore.NewJSONEncoder(cfg)

	var encoders []zapcore.Core
	var encoder zapcore.Encoder

	encoder = jsonEnc

	encoders = append(encoders,
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	)

	core := zapcore.NewTee(encoders...)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return log, nil
}
