package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/nash567/GoSentinel/pkg/logger/config"
	"github.com/nash567/GoSentinel/pkg/logger/model"
	"golang.org/x/exp/slog"
)

type SLogger struct {
	logger       *slog.Logger
	programLevel *slog.LevelVar
}

func NewSLogger(config *config.Config) *SLogger {
	opts := &slog.HandlerOptions{
		Level:       config.Level.SLogLevel(),
		ReplaceAttr: model.ReplaceAttributes,
	}
	log := &SLogger{
		logger:       slog.New(slog.NewJSONHandler(os.Stdout, opts)),
		programLevel: new(slog.LevelVar),
	}
	slog.SetDefault(log.logger)
	return log

}

func (l *SLogger) Trace(msg string, args ...any) {
	l.logger.Log(context.Background(), model.LevelTrace.SLogLevel(), msg, args...)
}

func (l *SLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args)
}
func (l *SLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *SLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SLogger) Fatal(msg string, args ...any) {
	l.logger.Log(context.Background(), model.LevelFatal.SLogLevel(), msg, args...)
	panic(msg)
}

func (l *SLogger) WithField(key string, value any) model.Logger {
	return &SLogger{
		logger:       l.logger.With(key, value),
		programLevel: l.programLevel,
	}
}
func (l *SLogger) WithFields(field model.Fields) model.Logger {
	attributes := make([]interface{}, 0)
	for key, value := range field {
		attributes = append(attributes, key, value)
	}
	return &SLogger{
		logger:       l.logger.With(attributes...),
		programLevel: l.programLevel,
	}
}

// TODO: Implement
//
//	func (l *SLogger) WithFields(key string, value any) {
//		l.logger.With(key, value)
//	}

func (l *SLogger) SetLevel(level model.Level) {
	l.programLevel.Set(level.SLogLevel())
}
