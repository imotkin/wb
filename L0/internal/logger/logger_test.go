package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

type logRecord struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
	Int     int    `json:"int"`
	String  string `json:"string"`
}

var logFunc = func(l Logger) func(msg string, args ...any) {
	switch l.Level() {
	case LevelDebug:
		return l.Debug
	case LevelInfo:
		return l.Info
	case LevelWarn:
		return l.Warn
	default:
		return nil
	}
}

func TestLoggerText(t *testing.T) {
	cases := []struct {
		level    Level
		msg      string
		args     []any
		expected []string
	}{
		{
			LevelDebug,
			"Hello, World!",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			[]string{"DEBUG", `msg="Hello, World!"`, "int=123", "string=text"},
		},
		{
			LevelInfo,
			"Log message",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			[]string{"INFO", `msg="Log message"`, "int=123", "string=text"},
		},
		{
			LevelWarn,
			"message 123",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			[]string{"WARN", `msg="message 123"`, "int=123", "string=text"},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			var (
				buf bytes.Buffer
				log = New(FormatText, tt.level, &buf)
			)

			logFunc(log)(tt.msg, tt.args...)

			for _, s := range tt.expected {
				require.Contains(t, buf.String(), s)
			}
		})
	}
}

func TestLoggerJSON(t *testing.T) {
	cases := []struct {
		level    Level
		msg      string
		args     []any
		expected logRecord
	}{
		{
			LevelDebug,
			"Hello, World!",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			logRecord{
				Level:   "DEBUG",
				Message: "Hello, World!",
				Int:     123,
				String:  "text",
			},
		},
		{
			LevelInfo,
			"...",
			[]any{slog.Int("int", 0), slog.String("string", "")},
			logRecord{
				Level:   "INFO",
				Message: "...",
				Int:     0,
				String:  "",
			},
		},
		{
			LevelWarn,
			"Warning...",
			[]any{slog.Int("int", -1), slog.String("string", "///")},
			logRecord{
				Level:   "WARN",
				Message: "Warning...",
				Int:     -1,
				String:  "///",
			},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			var (
				buf bytes.Buffer
				got logRecord
			)

			l := New(FormatJSON, tt.level, &buf)

			logFunc(l)(tt.msg, tt.args...)

			err := json.Unmarshal(buf.Bytes(), &got)
			require.NoError(t, err)

			require.Equal(t, tt.expected, got)
		})
	}
}

func TestLoggerTextError(t *testing.T) {
	cases := []struct {
		err      error
		msg      string
		args     []any
		expected []string
	}{
		{
			errors.New("My Error"),
			"something has failed...",
			[]any{},
			[]string{"level=ERROR", `msg="something has failed..."`, `error="My Error"`},
		},
		{
			errors.New("another error!"),
			"fail",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			[]string{"level=ERROR", `msg=fail`, `error="another error!"`, "int=123", "string=text"},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			var (
				buf = new(bytes.Buffer)
				l   = New(FormatText, LevelError, buf)
			)

			l.Error(tt.err, tt.msg, tt.args...)

			for _, s := range tt.expected {
				require.Contains(t, buf.String(), s)
			}
		})
	}
}

func TestLoggerJSONError(t *testing.T) {
	cases := []struct {
		err      error
		msg      string
		args     []any
		expected []string
	}{
		{
			errors.New("My Error"),
			"something has failed...",
			[]any{},
			[]string{"level=ERROR", `error="My Error"`},
		},
		{
			errors.New("another error!"),
			"something has failed...",
			[]any{slog.Int("int", 123), slog.String("string", "text")},
			[]string{"level=ERROR", `error="another error!"`, "int=123", "string=text"},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			var (
				buf = new(bytes.Buffer)
				l   = New(FormatText, LevelError, buf)
			)

			l.Error(tt.err, tt.msg, tt.args...)

			for _, s := range tt.expected {
				require.Contains(t, buf.String(), s)
			}
		})
	}
}
