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

func newHeatmapTestServer(t *testing.T) (*srv.Server, *mocks.MockArtistService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockArtistService(ctrl)
	return srv.New(svc, zap.NewNop()), svc
}

func TestGetHeatmapHandler_HappyPath(t *testing.T) {
	s, svc := newHeatmapTestServer(t)

	want := []lib.HeatmapCell{
		{Dow: 1, Hour: 8, Plays: 45},
		{Dow: 5, Hour: 22, Plays: 91},
	}
	svc.EXPECT().
		GetHeatmap(gomock.Any(), "artist-1", gomock.Any(), gomock.Any()).
		Return(want, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/heatmap?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Heatmap []lib.HeatmapCell `json:"heatmap"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Heatmap) != 2 || resp.Heatmap[0].Plays != 45 {
		t.Errorf("unexpected heatmap: %+v", resp.Heatmap)
	}
}

func TestGetHeatmapHandler_MissingArtistID(t *testing.T) {
	s, _ := newHeatmapTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/heatmap?start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetHeatmapHandler_MalformedStart(t *testing.T) {
	s, _ := newHeatmapTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/heatmap?artist_id=artist-1&start=notadate&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetHeatmapHandler_ServiceError(t *testing.T) {
	s, svc := newHeatmapTestServer(t)
	svc.EXPECT().
		GetHeatmap(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet,
		"/heatmap?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status %d, want 500", w.Code)
	}
}

func TestGetHeatmapHandler_EmptyResult(t *testing.T) {
	s, svc := newHeatmapTestServer(t)
	svc.EXPECT().
		GetHeatmap(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]lib.HeatmapCell{}, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/heatmap?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Heatmap []lib.HeatmapCell `json:"heatmap"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Heatmap == nil || len(resp.Heatmap) != 0 {
		t.Errorf("expected empty array, got %+v", resp.Heatmap)
	}
}
