package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type (
	Format string
	Level  string
	Output string
)

func (f Format) Validate() error {
	switch f {
	case FormatJSON, FormatText:
		return nil
	default:
		return validation.ErrInInvalid
	}
}

func (l Level) Validate() error {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return nil
	default:
		return validation.ErrInInvalid
	}
}

func (o Output) Validate() error {
	switch {
	case strings.TrimSpace(string(o)) != "":
		return nil
	default:
		return validation.ErrRequired
	}
}

func ParseOutput(o Output) (io.Writer, error) {
	switch o {
	case "stdout", "out":
		return os.Stdout, nil
	case "stderr", "err":
		return os.Stderr, nil
	default:
		return os.Create(string(o))
	}
}

func ParseLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func New(format Format, level Level, w io.Writer) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: ParseLevel(level),
	}

	switch format {
	case FormatJSON:
		handler = slog.NewJSONHandler(w, opts)
	case FormatText:
		handler = slog.NewTextHandler(w, opts)
	}

	return &logger{
		l:     slog.New(handler),
		level: level,
	}
}

func NewNoOp() Logger {
	return New(FormatJSON, LevelInfo, io.Discard)
}
