package logger

import (
	"log/slog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(err error, msg string, args ...any)
	With(args ...any) Logger
}

type logger struct {
	l *slog.Logger
}

func (l *logger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.l.Warn(msg, args...)
}

func (l *logger) Error(err error, msg string, args ...any) {
	attr := slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
	l.l.Error(msg, append(args, attr)...)
}

func (l *logger) With(args ...any) Logger {
	return &logger{l.l.With(args...)}
}
