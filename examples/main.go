package main

import (
	"os"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slog"

	"github.com/jkratz55/zeroslog"
)

func main() {

	slogger := slog.New(slog.NewJSONHandler(os.Stderr))
	slogger.Info("Hello")

	zLogger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		CallerWithSkipFrameCount(5).
		Stack().
		Logger()
	zeroSlogger := slog.New(zeroslog.NewHandler(zLogger))

	zeroSlogger.Error("Hey this is a test",
		slog.String("serviceId", "34320feihsfhesfhesi"))
}
