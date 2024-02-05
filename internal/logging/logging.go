package logging

import (
	"io"
	"log/slog"
	"os"
	"path"
)

var appLogLevel slog.LevelVar = slog.LevelVar{}

type LoggingOptions struct {
	LogFilePath string
	LogLevel    slog.Level
}

func init() {
	configure(slog.LevelInfo, os.Stdout)
}

func Configure(options LoggingOptions) {
	var writer io.Writer

	if options.LogFilePath != "" {
		dir := path.Dir(options.LogFilePath)
		os.MkdirAll(dir, 0755)
		f, err := os.OpenFile(options.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		writer = io.MultiWriter(os.Stdout, f)
	} else {
		writer = os.Stdout
	}

	configure(options.LogLevel, writer)
	slog.Default().Debug("Configured logger")
}

func configure(level slog.Level, writer io.Writer) {
	appLogLevel.Set(level)
	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: &appLogLevel,
	}))

	slog.SetDefault(logger)
}

func SetLogLevel(level slog.Level) {
	appLogLevel.Set(level)
}
