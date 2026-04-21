package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imotkin/L2/18/internal/calendar"
	"github.com/matryer/is"
	"go.uber.org/mock/gomock"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func buildRequest(method string, v any) (*http.Request, error) {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}

	return httptest.NewRequest(method, "localhost:8080", buf), nil
}

func TestHandlerCreateEvent(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		is   = is.New(t)
		w    = httptest.NewRecorder()
		s    = calendar.NewMockService(ctrl)
		h    = NewHandler(logger, s)

		r = calendar.CreateRequest{
			UserID: uuid.New(),
			Date:   time.Now().Round(0),
			Text:   "Hello!",
		}

		event = r.ToEvent()
	)

	req, err := buildRequest(http.MethodPost, r)
	is.NoErr(err)

	s.EXPECT().
		CreateEvent(context.Background(), r).
		Return(event, nil)

	h.CreateEvent(w, req)

	var created calendar.Event
	err = json.NewDecoder(w.Body).Decode(&created)
	is.NoErr(err)

	is.Equal(w.Result().StatusCode, http.StatusOK)
	is.Equal(event, created)
}

func TestHandlerDeleteEvent(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		is   = is.New(t)
		w    = httptest.NewRecorder()
		s    = calendar.NewMockService(ctrl)
		h    = NewHandler(logger, s)
		buf  = new(bytes.Buffer)
	)

	createReq := calendar.CreateRequest{
		UserID: uuid.New(),
		Date:   time.Now().Round(0),
		Text:   "Hello!",
	}

	event := createReq.ToEvent()

	err := json.NewEncoder(buf).Encode(createReq)
	is.NoErr(err)

	req := httptest.NewRequest(
		http.MethodPost, "localhost:8080", buf,
	)

	s.EXPECT().
		CreateEvent(context.Background(), createReq).
		Return(event, nil)

	h.CreateEvent(w, req)

	var created calendar.Event
	err = json.NewDecoder(w.Body).Decode(&created)
	is.NoErr(err)

	is.Equal(w.Result().StatusCode, http.StatusOK)
	is.Equal(event, created)

	delReq := calendar.DeleteRequest{
		UserID:  event.UserID,
		EventID: event.ID,
	}

	req, err = buildRequest(http.MethodPost, delReq)

	s.EXPECT().
		DeleteEvent(context.Background(), delReq).
		Return(event, nil)

	h.DeleteEvent(w, req)

	var deleted calendar.Event
	err = json.NewDecoder(w.Body).Decode(&deleted)
	is.NoErr(err)

	is.Equal(w.Result().StatusCode, http.StatusOK)
	is.Equal(event, created)
}
