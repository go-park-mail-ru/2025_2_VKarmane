package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type loggerKeyType string

const loggerKey loggerKeyType = "logger"

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return nil
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type SlogLogger struct {
	logger *slog.Logger
	closer io.Closer
}

func NewSlogLogger() Logger {
	return &SlogLogger{
		logger: slog.Default(),
	}
}

func NewSlogLoggerWithLevel(level slog.Level) Logger {
	var handler slog.Handler
	switch level {
	case slog.LevelDebug:
		handler = slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelDebug})
	case slog.LevelInfo:
		handler = slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelInfo})
	case slog.LevelWarn:
		handler = slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelWarn})
	case slog.LevelError:
		handler = slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError})
	default:
		handler = slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func NewSlogLoggerWithFile(filename string, level slog.Level) (Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})

	return &SlogLogger{
		logger: slog.New(handler),
		closer: file,
	}, nil
}

func NewSlogLoggerWithFileAndConsole(filename string, level slog.Level) (Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := &MultiWriter{
		file:    file,
		console: os.Stdout,
	}

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})

	return &SlogLogger{
		logger: slog.New(handler),
		closer: file,
	}, nil
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *SlogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *SlogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *SlogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}


func (l *SlogLogger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

type MultiWriter struct {
	file    *os.File
	console *os.File
}

func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	if _, err := mw.file.Write(p); err != nil {
		return 0, err
	}

	if _, err := mw.console.Write(p); err != nil {
		return 0, err
	}

	return len(p), nil
}
