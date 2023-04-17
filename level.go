package zeroslog

import (
	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"
)

type Level = slog.Level

const (
	LevelTrace Level = slog.Level(-5)
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
	LevelPanic Level = slog.Level(9)
	LevelFatal Level = slog.Level(10)
)

func mapSlogLevelToZerolog(lvl slog.Level) zerolog.Level {
	switch lvl {
	case LevelTrace:
		return zerolog.TraceLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelPanic:
		return zerolog.PanicLevel
	case LevelFatal:
		return zerolog.PanicLevel
	default:
		// If we fall into the default block we have no idea what level the record
		// is trying to be logged at, so we just set it to NoLevel and move on with
		// life.
		return zerolog.NoLevel
	}
}

func mapZerologLevelToSlog(lvl zerolog.Level) slog.Level {
	switch lvl {
	case zerolog.TraceLevel:
		return LevelTrace
	case zerolog.DebugLevel:
		return LevelDebug
	case zerolog.InfoLevel:
		return LevelInfo
	case zerolog.WarnLevel:
		return LevelWarn
	case zerolog.ErrorLevel:
		return LevelError
	case zerolog.PanicLevel:
		return LevelPanic
	case zerolog.FatalLevel:
		return LevelFatal
	default:
		return LevelTrace
	}
}
