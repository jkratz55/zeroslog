package zeroslog

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"
)

type Handler struct {
	logger    zerolog.Logger
	group     *zerolog.Event
	groupName string
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

	// Map slog level to the zerolog level
	lvl := mapSlogLevelToZerolog(record.Level)

	// Create an event for the zerolog mapped level from the slog level. If the
	// level is ERROR, PANIC, or FATAL append the stacktrace to the event.
	event := h.logger.WithLevel(lvl)
	if lvl == zerolog.ErrorLevel || lvl == zerolog.PanicLevel || lvl == zerolog.FatalLevel {
		event.Str("stacktrace", string(debug.Stack()))
	}

	if h.group != nil {
		dict := zerolog.Dict()
		record.Attrs(func(attr slog.Attr) {
			appendAttr(dict, attr)
		})
		event.Dict(h.groupName, dict)
		event.Msg(record.Message)
	} else {
		record.Attrs(func(attr slog.Attr) {
			event = appendAttr(event, attr)
		})
		event.Msg(record.Message)
	}
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
	if strings.TrimSpace(name) == "" {
		return h
	}
	if h.group == nil {
		return &Handler{
			logger:    h.logger,
			group:     zerolog.Dict(),
			groupName: name,
		}
	}
	return &Handler{
		logger:    h.logger,
		group:     h.group.Dict(h.groupName, h.group),
		groupName: name,
	}
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
