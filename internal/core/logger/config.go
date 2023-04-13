package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Conf struct {
	Level           string `mapstructure:"level" yaml:"level"`
	Encoding        string `mapstructure:"encoding" yaml:"encoding"` // json or console
	Directory       string `mapstructure:"directory"  yaml:"directory"`
	LogLevelEncoder string `mapstructure:"log-level-encoder" yaml:"log-level-encoder"`
	MaxAge          int    `mapstructure:"max-age" yaml:"max-age"`
	ShowLine        bool   `mapstructure:"show-line" yaml:"show-line"`
	LogInConsole    bool   `mapstructure:"log-in-console" yaml:"log-in-console"`
	Compress        bool   `mapstructure:"compress" yaml:"compress"`
}

// ZapLogLevelEncoder string to zapcore.LevelEncoder
func (z *Conf) ZapLogLevelEncoder() zapcore.LevelEncoder {
	switch {
	case z.LogLevelEncoder == "LowercaseLevelEncoder": // e.g. info
		return zapcore.LowercaseLevelEncoder
	case z.LogLevelEncoder == "LowercaseColorLevelEncoder": // e.g.  info (colorful)
		return zapcore.LowercaseColorLevelEncoder
	case z.LogLevelEncoder == "CapitalLevelEncoder": // e.g.  INFO
		return zapcore.CapitalLevelEncoder
	case z.LogLevelEncoder == "CapitalColorLevelEncoder": // e.g.  INFO (colorful)
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.CapitalColorLevelEncoder
	}
}

// TransformLevel  string to zapcore.Level
func (z *Conf) TransformLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
