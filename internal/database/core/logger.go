package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"my_db/internal/database/config"
)

func MustLoadLogger(conf config.Logging) *zap.Logger {
	dispatchLogLevels := map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
	}
	zapLogLevel, ok := dispatchLogLevels[conf.Level]
	if !ok {
		panic(fmt.Sprintf("unknown log level: %s", conf.Level))
	}

	loggerConf := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLogLevel),
		Encoding:         conf.Encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{conf.Output},
		ErrorOutputPaths: []string{conf.Output},
	}

	log, err := loggerConf.Build()
	if err != nil {
		panic(err)
	}
	return log
}
