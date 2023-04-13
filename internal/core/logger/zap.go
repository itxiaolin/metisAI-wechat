package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	globalConf   Conf
)

func init() {
	zapLogger := BuildZapLogger(Conf{
		Level:        "info",
		Encoding:     "console",
		Directory:    "./logs",
		MaxAge:       3,
		ShowLine:     true,
		LogInConsole: true,
		Compress:     true,
	})
	globalLogger = zapLogger
}

func BuildZapLogger(conf Conf) *zap.Logger {
	globalConf = conf

	cores := getZapCores()
	globalLogger = zap.New(zapcore.NewTee(cores...))

	if globalConf.ShowLine {
		globalLogger = globalLogger.WithOptions(zap.AddCaller())
	}
	return globalLogger
}

func SetZapLogger(logger *zap.Logger) {
	globalLogger = logger
}

func GetZapLogger() *zap.Logger {
	if globalLogger == nil {
		return zap.L()
	}
	return globalLogger
}

func getDefaultEncoderConf() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    globalConf.ZapLogLevelEncoder(),
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := globalConf.TransformLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, getEncoderCore(level, getLevelPriority(level)))
	}
	return cores
}

func getEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := getWriteSyncer(l.String())
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02-15:04:05.000"))
}

func getEncoder() zapcore.Encoder {
	if globalConf.Encoding == "json" {
		return zapcore.NewJSONEncoder(getDefaultEncoderConf())
	}
	return zapcore.NewConsoleEncoder(getDefaultEncoderConf())
}

// getWriteSyncer use lumberjack.Logger to write log file
func getWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(globalConf.Directory, level+".log"),
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     globalConf.MaxAge,
		Compress:   true,
	}
	fileWriter := zapcore.AddSync(lumberJackLogger)
	if globalConf.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), nil
	}
	return zapcore.AddSync(fileWriter), nil
}

func getLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool {
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool {
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool {
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool {
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	}
}
