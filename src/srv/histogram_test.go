package srv_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/el10savio/artistInsights/src/lib"
	"github.com/el10savio/artistInsights/src/srv/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/el10savio/artistInsights/src/srv"
)

func newHistogramTestServer(t *testing.T) (*srv.Server, *mocks.MockArtistService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockArtistService(ctrl)
	return srv.New(svc, zap.NewNop()), svc
}

func TestGetHistogramHandler_HappyPath(t *testing.T) {
	s, svc := newHistogramTestServer(t)

	want := []lib.HistogramBin{
		{Bin: "1–9", Count: 100},
		{Bin: "10–99", Count: 30},
		{Bin: "100–999", Count: 5},
		{Bin: "1000+", Count: 0},
	}
	svc.EXPECT().
		GetHistogram(gomock.Any(), "artist-1", gomock.Any(), gomock.Any(), "listener").
		Return(want, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/histogram?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&by=listener", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Histogram []lib.HistogramBin `json:"histogram"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Histogram) != 4 {
		t.Errorf("got %d bins, want 4", len(resp.Histogram))
	}
}

func TestGetHistogramHandler_MissingArtistID(t *testing.T) {
	s, _ := newHistogramTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/histogram?start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&by=listener", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetHistogramHandler_MalformedStart(t *testing.T) {
	s, _ := newHistogramTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/histogram?artist_id=artist-1&start=notadate&end=2022-12-31T23:59:59Z&by=listener", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetHistogramHandler_InvalidBy(t *testing.T) {
	s, _ := newHistogramTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/histogram?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&by=song", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetHistogramHandler_ServiceError(t *testing.T) {
	s, svc := newHistogramTestServer(t)
	svc.EXPECT().
		GetHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet,
		"/histogram?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&by=track", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status %d, want 500", w.Code)
	}
}

func TestGetHistogramHandler_EmptyResult(t *testing.T) {
	s, svc := newHistogramTestServer(t)
	svc.EXPECT().
		GetHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]lib.HistogramBin{}, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/histogram?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&by=listener", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Histogram []lib.HistogramBin `json:"histogram"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Histogram == nil || len(resp.Histogram) != 0 {
		t.Errorf("expected empty array, got %+v", resp.Histogram)
	}
}
