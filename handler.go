package zeroslog

import (
	"context"

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
	event := h.logger.WithLevel(lvl)

	record.Attrs(func(attr slog.Attr) {
		event.Str(attr.Key, attr.Value.String())
	})

	event.Msg(record.Message)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) WithGroup(name string) slog.Handler {
	// TODO implement me
	panic("implement me")
}
