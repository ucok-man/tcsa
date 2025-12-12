package tlog

import (
	"io"
	"os"
	"sync"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
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
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return NewLogger(logger), nil
}

func NewDevelopment() (*Logger, error) {
	logger, err := zap.NewDevelopment()
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

// shouldLog checks if the message should be logged based on level
func (l *Logger) shouldLog(lvl log.Lvl) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return lvl >= l.level
}

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
		msg, found := j["message"]
		if !found {
			l.withPrefix().Infow("", l.jsonToFields(j)...)
		}

		switch m := msg.(type) {
		case string:
			delete(j, "message")
			l.withPrefix().Infow(m, l.jsonToFields(j)...)
		default:
			l.withPrefix().Infow("", l.jsonToFields(j)...)
		}

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
		msg, found := j["message"]
		if !found {
			l.withPrefix().Debugw("", l.jsonToFields(j)...)
		}

		switch m := msg.(type) {
		case string:
			delete(j, "message")
			l.withPrefix().Debugw(m, l.jsonToFields(j)...)
		default:
			l.withPrefix().Debugw("", l.jsonToFields(j)...)
		}
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
		msg, found := j["message"]
		if !found {
			l.withPrefix().Infow("", l.jsonToFields(j)...)
		}

		switch m := msg.(type) {
		case string:
			delete(j, "message")
			l.withPrefix().Infow(m, l.jsonToFields(j)...)
		default:
			l.withPrefix().Infow("", l.jsonToFields(j)...)
		}
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
		msg, found := j["message"]
		if !found {
			l.withPrefix().Warnw("", l.jsonToFields(j)...)
		}

		switch m := msg.(type) {
		case string:
			delete(j, "message")
			l.withPrefix().Warnw(m, l.jsonToFields(j)...)
		default:
			l.withPrefix().Warnw("", l.jsonToFields(j)...)
		}
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
		msg, found := j["message"]
		if !found {
			l.withPrefix().Errorw("", l.jsonToFields(j)...)
		}

		switch m := msg.(type) {
		case string:
			delete(j, "message")
			l.withPrefix().Errorw(m, l.jsonToFields(j)...)
		default:
			l.withPrefix().Errorw("", l.jsonToFields(j)...)
		}
	}
}

func (l *Logger) Fatal(i ...interface{}) {
	l.withPrefix().Fatal(i...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.withPrefix().Fatalf(format, args...)
}

func (l *Logger) Fatalj(j log.JSON) {
	msg, found := j["message"]
	if !found {
		l.withPrefix().Fatalw("", l.jsonToFields(j)...)
	}

	switch m := msg.(type) {
	case string:
		delete(j, "message")
		l.withPrefix().Fatalw(m, l.jsonToFields(j)...)
	default:
		l.withPrefix().Fatalw("", l.jsonToFields(j)...)
	}
}

func (l *Logger) Panic(i ...interface{}) {
	l.withPrefix().Panic(i...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.withPrefix().Panicf(format, args...)
}

func (l *Logger) Panicj(j log.JSON) {
	msg, found := j["message"]
	if !found {
		l.withPrefix().Panicw("", l.jsonToFields(j)...)
	}

	switch m := msg.(type) {
	case string:
		delete(j, "message")
		l.withPrefix().Panicw(m, l.jsonToFields(j)...)
	default:
		l.withPrefix().Panicw("", l.jsonToFields(j)...)
	}
}

func (l *Logger) Sync() {
	l.logger.Sync()
	l.sugar.Sync()
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.withPrefix().Error(string(p))
	return len(p), nil
}
