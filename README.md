# zeroslog

Zeroslog is an implementation of the `slog.Handler` interface allowing Zerolog to be used as the backend logger for `slog`. 

_Notice: `slog` is still experimental and not officially part of the standard library. There may be changes to `slog` before it becomes finalized and part of the standard library. You may want to keep that in mind before using `slog` or this library._

_Warning: This library is experimental and not recommended for production use._

## Why Slog?

`slog` provides structured logging functionality from the Go team. Although experimental at this time, ideally it will become the defacto logging interface used by frameworks and libraries. Presently the logging ecosystem in Go is fragmented although there are a few excellent third party libraries such as Zap and Zerolog. Third party loggers can be used as Handlers for `slog` making it very versatile for library developers.

## Important Callouts 

* `WithGroup` on the `Handler` doesn't support nesting as in it won't allow a group within a group. The last one will win and be the group used.
* The `zerolog.Logger` needs to be configured with `CallerWithSkipFrameCount(5)` to report the correct location the log event was generated from. This is less than ideal and fragile.

# Example

```go
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
```