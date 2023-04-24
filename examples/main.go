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
	slogger.WithGroup("groupy").Info("Hello",
		slog.String("firstName", "Sir Elton"))

	slogger.WithGroup("group1").
		WithGroup("group2").
		WithGroup("group3").
		Info("Hello World",
			slog.String("firstName", "Sir Elton"))

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
			slog.String("browser", "CHROME"),
			slog.Group("sub",
				slog.String("name", "POWER"))))

	zeroSlogger.WithGroup("groupy").Info("Hello",
		slog.String("firstName", "Sir Elton"))

	zeroSlogger.WithGroup("groupy").
		WithGroup("group2").
		Info("Hello",
			slog.String("firstName", "Sir Elton"))
}
