package logging

import (
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"gin-example/pkg/setting"
)

var (
	Logger       *zap.Logger
	routerLogger *zap.Logger
	GormLogger   *zap.Logger
)

// Setup initialize the log instance
func Setup() {
	Logger = NewZapLogger(newAllFieldsZapCore())
	Logger.Info("initialization zap logger ok.", zap.String("LogLever", switchLogLevel().String()))

	routerLogger = NewZapLogger(newGinDebugZapCore())
	Logger.Info("initialization gin debug logger ok.", zap.String("LogLever", switchLogLevel().String()))

	GormLogger = NewZapWithoutCallerLogger(newGormZapCore())
	Logger.Info("initialization zap logger with caller ok.", zap.String("LogLever", switchLogLevel().String()))
}

func switchLogLevel() zapcore.Level {
	switch strings.ToLower(setting.ServerSetting.RunMode) {
	case "debug":
		return zapcore.DebugLevel
	case "release":
		return zapcore.InfoLevel
	}

	switch strings.ToLower(setting.LoggerSetting.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func switchLogOutput() zapcore.WriteSyncer {
	if setting.LoggerSetting.Stdout {
		return zapcore.AddSync(os.Stdout)
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   setting.LoggerSetting.FilePath, // 日志文件路径
		MaxSize:    128,                            // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                             // 日志文件最多保存多少个备份
		MaxAge:     7,                              // 文件最多保存多少天
		Compress:   true,                           // 是否压缩
	})
}

func ISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewZapLogger(core zapcore.Core) *zap.Logger {
	return zap.New(core, zap.AddCaller())
}

func NewZapWithoutCallerLogger(core zapcore.Core) *zap.Logger {
	return zap.New(core)
}

func newAllFieldsZapCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			NameKey:        "name",
			TimeKey:        "@timestamp",
			LevelKey:       "level",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder, //
			EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
			EncodeName:     zapcore.FullNameEncoder,
		}), // 编码器配置
		zapcore.NewMultiWriteSyncer(switchLogOutput()), // 打印到控制台或文件
		zap.NewAtomicLevelAt(switchLogLevel()),         // 日志级别
	)
}

func newGinDebugZapCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder, //
			EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
			EncodeName:     zapcore.FullNameEncoder,
		}), // 编码器配置
		zapcore.NewMultiWriteSyncer(switchLogOutput()), // 打印到控制台或文件
		zap.NewAtomicLevelAt(switchLogLevel()),         // 日志级别
	)
}

func newGormZapCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			NameKey:        "name",
			TimeKey:        "@timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder, //
			EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
			EncodeName:     zapcore.FullNameEncoder,
		}), // 编码器配置
		zapcore.NewMultiWriteSyncer(switchLogOutput()), // 打印到控制台或文件
		zap.NewAtomicLevelAt(switchLogLevel()),         // 日志级别
	)
}

// print gin route function
func GinDebugPrintRouteZapLoggerFunc(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	routerLogger.Debug("",
		zap.Strings("Gin Debug", []string{httpMethod, absolutePath, handlerName, strconv.Itoa(nuHandlers)}),
	)
}
