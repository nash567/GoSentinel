package model

import (
	"golang.org/x/exp/slog"
)

type Level int

const (
	LevelTrace Level = -8
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12
)

var LevelNames = map[slog.Leveler]string{
	slog.Level(LevelTrace): "TRACE",
	slog.Level(LevelFatal): "FATAL",
}

func ReplaceAttributes(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := LevelNames[level]
		if !exists {
			levelLabel = level.String()
		}
		a.Value = slog.StringValue(levelLabel)

	}
	return a

}

func (l Level) SLogLevel() slog.Level {
	switch l {
	case LevelTrace:
		return slog.Level(LevelTrace)
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	case LevelFatal:
		return slog.Level(LevelFatal)
	default:
		return slog.LevelInfo
	}
}
