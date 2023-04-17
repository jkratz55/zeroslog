package main

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/exp/slog"

	"github.com/jkratz55/zeroslog"
)

func main() {

	slogger := slog.New(slog.NewJSONHandler(os.Stderr))
	slogger.Info("Hello")

	zerolog.ErrorStackFieldName = "stack"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zLogger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		CallerWithSkipFrameCount(5).
		Stack().
		Logger()
	zeroSlogger := slog.New(zeroslog.NewHandler(zLogger))

	zeroSlogger.Info("Hey this is a test",
		slog.String("serviceId", "34320feihsfhesfhesi"),
		slog.String("userId", "jkratz"),
		slog.Group("http",
			slog.String("method", "GET"),
			slog.String("browser", "CHROME")))
}

func something() {
	something2()
}

func something2() {
	something3()
}

func something3() {
	something4()
}

func something4() {
	os.Stderr.Write(debug.Stack())
}
