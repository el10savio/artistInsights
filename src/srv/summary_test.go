package srv_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/el10savio/artistInsights/src/lib"
	"github.com/el10savio/artistInsights/src/srv/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/el10savio/artistInsights/src/srv"
)

func newTestServer(t *testing.T) (*srv.Server, *mocks.MockArtistService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockArtistService(ctrl)
	s := srv.New(svc, zap.NewNop())
	return s, svc
}

func TestGetSummaryHandler_HappyPath(t *testing.T) {
	s, svc := newTestServer(t)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	want := &lib.Summary{TotalPlays: 50, TotalUniqueListeners: 10}

	svc.EXPECT().
		GetSummary(gomock.Any(), "artist-1", start, end).
		Return(want, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/summary?artist_id=artist-1&start=2020-01-01T00:00:00Z&end=2021-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var got lib.Summary
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != *want {
		t.Errorf("got %+v, want %+v", got, *want)
	}
}

func TestGetSummaryHandler_MissingArtistID(t *testing.T) {
	s, _ := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet,
		"/summary?start=2020-01-01T00:00:00Z&end=2021-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSummaryHandler_MissingStart(t *testing.T) {
	s, _ := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet,
		"/summary?artist_id=artist-1&end=2021-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSummaryHandler_MalformedStart(t *testing.T) {
	s, _ := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet,
		"/summary?artist_id=artist-1&start=notadate&end=2021-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSummaryHandler_ServiceError(t *testing.T) {
	s, svc := newTestServer(t)

	svc.EXPECT().
		GetSummary(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet,
		"/summary?artist_id=artist-1&start=2020-01-01T00:00:00Z&end=2021-01-01T00:00:00Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status %d, want 500", w.Code)
	}
}
