package zeroslog

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"
)

type Handler struct {
	logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	loggerLevel := mapZerologLevelToSlog(h.logger.GetLevel())
	return lvl >= loggerLevel
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	lvl := mapSlogLevelToZerolog(record.Level)
	var event *zerolog.Event
	switch lvl {
	case zerolog.TraceLevel:
		event = h.logger.Trace()
	case zerolog.DebugLevel:
		event = h.logger.Debug()
	case zerolog.InfoLevel:
		event = h.logger.Info()
	case zerolog.WarnLevel:
		event = h.logger.Warn()
	case zerolog.ErrorLevel:
		event = h.logger.Error().Str("stacktrace", string(debug.Stack()))
	case zerolog.PanicLevel:
		event = h.logger.Panic().Str("stacktrace", string(debug.Stack()))
	case zerolog.FatalLevel:
		event = h.logger.Fatal().Str("stacktrace", string(debug.Stack()))
	default:
		event = h.logger.Log()
	}

	record.Attrs(func(attr slog.Attr) {
		event = appendAttr(event, attr)
	})

	event.Msg(record.Message)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	ctx := h.logger.With()
	for _, attr := range attrs {
		ctx = appendCtx(ctx, attr)
	}
	return &Handler{
		logger: ctx.Logger(),
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	// TODO implement me
	panic("implement me")
}

func appendAttr(evt *zerolog.Event, attr slog.Attr) *zerolog.Event {

	// Depending on the kind we can simply handle the type by called a method on
	// the slog.Value type to get the real value.
	if attr.Value.Kind() != slog.KindGroup {
		switch attr.Value.Kind() {
		case slog.KindBool:
			return evt.Bool(attr.Key, attr.Value.Bool())
		case slog.KindDuration:
			return evt.Dur(attr.Key, attr.Value.Duration())
		case slog.KindFloat64:
			return evt.Float64(attr.Key, attr.Value.Float64())
		case slog.KindInt64:
			return evt.Int64(attr.Key, attr.Value.Int64())
		case slog.KindString:
			return evt.Str(attr.Key, attr.Value.String())
		case slog.KindTime:
			return evt.Time(attr.Key, attr.Value.Time())
		case slog.KindUint64:
			return evt.Uint64(attr.Key, attr.Value.Uint64())
		}

		// If the kind didn't match the above block we need to check the type of raw
		// value.
		switch val := attr.Value.Any().(type) {
		case error:
			return evt.Err(val)
		case []error:
			return evt.Stack().Errs(attr.Key, val)
		case fmt.Stringer:
			return evt.Stringer(attr.Key, val)
		default:
			return evt.Any(attr.Key, val)
		}
	}

	if len(attr.Value.Group()) > 0 {
		dict := zerolog.Dict()
		for _, groupAttr := range attr.Value.Group() {
			appendAttr(dict, groupAttr)
		}
		return evt.Dict(attr.Key, dict)
	}

	return nil
}

func appendCtx(ctx zerolog.Context, attr slog.Attr) zerolog.Context {
	// Depending on the kind we can simply handle the type by called a method on
	// the slog.Value type to get the real value.
	switch attr.Value.Kind() {
	case slog.KindBool:
		return ctx.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return ctx.Dur(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return ctx.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return ctx.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return ctx.Str(attr.Key, attr.Value.String())
	case slog.KindTime:
		return ctx.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return ctx.Uint64(attr.Key, attr.Value.Uint64())
	}

	// If the kind didn't match the above block we need to check the type of raw
	// value.
	switch val := attr.Value.Any().(type) {
	case error:
		return ctx.Err(val)
	case []error:
		return ctx.Stack().Errs(attr.Key, val)
	case fmt.Stringer:
		return ctx.Stringer(attr.Key, val)
	default:
		// no-op, zerolog.Context doesn't support Any
		return ctx
	}
}
