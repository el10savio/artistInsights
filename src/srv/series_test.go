package srv_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/el10savio96/artistInsights/src/lib"
	"github.com/el10savio96/artistInsights/src/srv/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/el10savio96/artistInsights/src/srv"
)

func newSeriesTestServer(t *testing.T) (*srv.Server, *mocks.MockArtistService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockArtistService(ctrl)
	return srv.New(svc, zap.NewNop()), svc
}

func TestGetSeriesHandler_HappyPath(t *testing.T) {
	s, svc := newSeriesTestServer(t)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)
	want := []lib.SeriesPoint{
		{Bucket: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC), TotalPlays: 2, TotalUniqueListeners: 1},
	}

	svc.EXPECT().
		GetSeries(gomock.Any(), "artist-1", start, end, "month").
		Return(want, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/series?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&granularity=month", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Series []lib.SeriesPoint `json:"series"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Series) != 1 || resp.Series[0].TotalPlays != 2 {
		t.Errorf("unexpected series: %+v", resp.Series)
	}
}

func TestGetSeriesHandler_MissingArtistID(t *testing.T) {
	s, _ := newSeriesTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/series?start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&granularity=day", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSeriesHandler_MalformedStart(t *testing.T) {
	s, _ := newSeriesTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/series?artist_id=artist-1&start=notadate&end=2022-12-31T23:59:59Z&granularity=day", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSeriesHandler_InvalidGranularity(t *testing.T) {
	s, _ := newSeriesTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/series?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&granularity=minute", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetSeriesHandler_ServiceError(t *testing.T) {
	s, svc := newSeriesTestServer(t)
	svc.EXPECT().
		GetSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet,
		"/series?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&granularity=day", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status %d, want 500", w.Code)
	}
}

func TestGetSeriesHandler_EmptyResult(t *testing.T) {
	s, svc := newSeriesTestServer(t)
	svc.EXPECT().
		GetSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]lib.SeriesPoint{}, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/series?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z&granularity=day", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Series []lib.SeriesPoint `json:"series"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Series == nil || len(resp.Series) != 0 {
		t.Errorf("expected empty array, got %+v", resp.Series)
	}
}
