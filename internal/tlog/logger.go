package tlog

import (
	"io"
	"os"
	"sync"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type JSON = log.JSON

type SugaredLogger = zap.SugaredLogger

type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	output io.Writer
	level  log.Lvl
	prefix string
	mu     sync.RWMutex
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
		level:  log.INFO,
		output: os.Stdout,
	}
}

func NewProduction() (*Logger, error) {
	logger, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return NewLogger(logger), nil
}

func NewDevelopment() (*Logger, error) {
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return NewLogger(logger), nil
}

func Must(logger *Logger, err error) *Logger {
	if err != nil {
		panic(err)
	}
	return logger
}

func (l *Logger) Output() io.Writer {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.output
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writeSyncer := zapcore.AddSync(w)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.sugar = l.logger.Sugar()
}

func (l *Logger) Prefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.prefix
}

func (l *Logger) SetPrefix(p string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = p
}

func (l *Logger) Level() log.Lvl {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

func (l *Logger) SetLevel(v log.Lvl) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = v
}

func (l *Logger) SetHeader(h string) {}

func (l *Logger) WithSkipCaller(skip int) *Logger {
	logger := l.logger.WithOptions(zap.AddCallerSkip(skip))
	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
		level:  log.INFO,
		output: l.output,
	}

}

func (l *Logger) Print(i ...any) {
	if l.shouldLog(log.INFO) {
		l.withPrefix().Info(i)
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	if l.shouldLog(log.INFO) {
		l.withPrefix().Infof(format, args...)
	}
}

func (l *Logger) Printj(j log.JSON) {
	if l.shouldLog(log.INFO) {
		l.logJSON(j, l.withPrefix().Infow)
	}
}

func (l *Logger) Debug(i ...interface{}) {
	if l.shouldLog(log.DEBUG) {
		l.withPrefix().Debug(i...)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.shouldLog(log.DEBUG) {
		l.withPrefix().Debugf(format, args...)
	}
}

func (l *Logger) Debugj(j log.JSON) {
	if l.shouldLog(log.DEBUG) {
		l.logJSON(j, l.withPrefix().Debugw)
	}
}

func (l *Logger) Info(i ...interface{}) {
	if l.shouldLog(log.INFO) {
		l.withPrefix().Info(i...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.shouldLog(log.INFO) {
		l.withPrefix().Infof(format, args...)
	}
}

func (l *Logger) Infoj(j log.JSON) {
	if l.shouldLog(log.INFO) {
		l.logJSON(j, l.withPrefix().Infow)
	}
}

func (l *Logger) Warn(i ...interface{}) {
	if l.shouldLog(log.WARN) {
		l.withPrefix().Warn(i...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.shouldLog(log.WARN) {
		l.withPrefix().Warnf(format, args...)
	}
}

func (l *Logger) Warnj(j log.JSON) {
	if l.shouldLog(log.WARN) {
		l.logJSON(j, l.withPrefix().Warnw)
	}
}

func (l *Logger) Error(i ...interface{}) {
	if l.shouldLog(log.ERROR) {
		l.withPrefix().Error(i...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.shouldLog(log.ERROR) {
		l.withPrefix().Errorf(format, args...)
	}
}

func (l *Logger) Errorj(j log.JSON) {
	if l.shouldLog(log.ERROR) {
		l.logJSON(j, l.withPrefix().Errorw)
	}
}

func (l *Logger) Fatal(i ...interface{}) {
	l.withPrefix().Fatal(i...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.withPrefix().Fatalf(format, args...)
}

func (l *Logger) Fatalj(j log.JSON) {
	l.logJSON(j, l.withPrefix().Fatalw)
}

func (l *Logger) Panic(i ...interface{}) {
	l.withPrefix().Panic(i...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.withPrefix().Panicf(format, args...)
}

func (l *Logger) Panicj(j log.JSON) {
	l.logJSON(j, l.withPrefix().Panicw)
}

func (l *Logger) Sync() {
	l.logger.Sync()
	l.sugar.Sync()
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.withPrefix().Error(string(p))
	return len(p), nil
}

// shouldLog checks if the message should be logged based on level
func (l *Logger) shouldLog(lvl log.Lvl) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return lvl >= l.level
}

// withPrefix add prefix field to logger
func (l *Logger) withPrefix() *SugaredLogger {
	l.mu.RLock()
	prefix := l.prefix
	l.mu.RUnlock()

	if prefix != "" {
		return l.sugar.With("prefix", prefix)
	}
	return l.sugar
}

// jsonToFields converts log.JSON map to zap fields
func (l *Logger) jsonToFields(j log.JSON) []any {
	fields := make([]any, 0)
	for k, v := range j {
		fields = append(fields, k, v)
	}
	return fields
}

// logJSON handles the common logic for logging JSON with message extraction
func (l *Logger) logJSON(j log.JSON, logFunc func(string, ...any)) {
	msg, found := j["message"]
	if !found {
		logFunc("", l.jsonToFields(j)...)
		return
	}

	switch m := msg.(type) {
	case string:
		delete(j, "message")
		logFunc(m, l.jsonToFields(j)...)
	default:
		logFunc("", l.jsonToFields(j)...)
	}
}
