package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/imotkin/L0/internal/logger"
)

func TestErrorResponse(t *testing.T) {
	h := New(logger.NewNoOp(), nil, nil)

	cases := []struct {
		message string
		code    int
		err     error
	}{
		{
			message: "Server Error",
			code:    500,
			err:     errors.New("Server Error"),
		},
		{
			message: "NOT FOUND",
			code:    404,
			err:     errors.New("Not Found (404)"),
		},
		{
			message: "Bad Request",
			code:    400,
			err:     errors.New("oops..."),
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			r := httptest.NewRecorder()

			h.error(r, tt.message, tt.code, tt.err)

			require.Equal(t, tt.code, r.Code)

			var msg ErrorMessage

			err := json.NewDecoder(r.Body).Decode(&msg)
			require.NoError(t, err)

			require.Equal(t, tt.message, msg.Message)
			require.Equal(t, tt.code, msg.StatusCode)
			require.Equal(t, http.StatusText(tt.code), msg.StatusMessage)
		})
	}
}

func TestResponse(t *testing.T) {
	h := New(logger.NewNoOp(), nil, nil)

	cases := []struct {
		body     any
		code     int
		response string
		err      error
	}{
		{
			body: struct {
				Message string `json:"message"`
			}{"Hello, World!"},
			code:     200, // == http.StatusOK
			response: `{"message":"Hello, World!"}`,
		},
		{
			body:     "123",
			code:     201, // == http.StatusCreated
			response: `"123"`,
		},
		{
			body:     "text",
			code:     202, // == http.StatusAccepted
			response: `"text"`,
		},

		{
			body:     "",
			code:     204, // == http.StatusNoContent
			response: `""`,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			r := httptest.NewRecorder()

			h.response(r, tt.body, tt.code)

			require.Equal(t, tt.code, r.Code)
			require.Equal(t, tt.response, strings.TrimRight(r.Body.String(), "\n"))
		})
	}
}
