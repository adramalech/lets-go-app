package logger

import (
    "log"
    "os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
    sugaredLogger *zap.SugaredLogger
}

func newZapLogger(config Configuration) (Logger, error) {
    debugWriter := zapcore.Lock(os.Stdout)
    errorWriter := zapcore.Lock(os.Stderr)
    encoder := getEncoder(config.UseJSONFormat)

    highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= zapcore.ErrorLevel
    })

    lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl < zapcore.ErrorLevel
    })

    debugCore := zapcore.NewCore(encoder, debugWriter, lowPriority)

    errorCore := zapcore.NewCore(encoder, errorWriter, highPriority)

    combinedCores := zapcore.NewTee(debugCore, errorCore)

    logger := zap.New(combinedCores, zap.AddCallerSkip(2), zap.AddCaller()).Sugar()

    zLogger := &zapLogger{sugaredLogger: logger}

    return zLogger, nil
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *zapLogger) Close() {
    l.sugaredLogger.Sync()
}

func (l *zapLogger) GetStdLogger() *log.Logger {
    return zap.NewStdLog(l.sugaredLogger.Desugar())
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	
    for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	
    newLogger := l.sugaredLogger.With(f...)
	
    return &zapLogger{newLogger}
}

func getEncoder(isJSON bool) zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
    if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	
    return zapcore.NewConsoleEncoder(encoderConfig)
}
