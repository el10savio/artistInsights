package srv_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/el10savio96/artistInsights/src/lib"
	"github.com/el10savio96/artistInsights/src/srv/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/el10savio96/artistInsights/src/srv"
)

func newTopTracksTestServer(t *testing.T) (*srv.Server, *mocks.MockArtistService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockArtistService(ctrl)
	return srv.New(svc, zap.NewNop()), svc
}

func TestGetTopTracksHandler_HappyPath(t *testing.T) {
	s, svc := newTopTracksTestServer(t)

	want := []lib.TrackPlay{
		{TrackID: "track-1", TrackName: "Song One", Plays: 100},
		{TrackID: "track-2", TrackName: "Song Two", Plays: 50},
	}
	svc.EXPECT().
		GetTopTracks(gomock.Any(), "artist-1", gomock.Any(), gomock.Any()).
		Return(want, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/top-tracks?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Tracks []lib.TrackPlay `json:"tracks"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Tracks) != 2 || resp.Tracks[0].Plays != 100 {
		t.Errorf("unexpected tracks: %+v", resp.Tracks)
	}
}

func TestGetTopTracksHandler_MissingArtistID(t *testing.T) {
	s, _ := newTopTracksTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/top-tracks?start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetTopTracksHandler_MalformedStart(t *testing.T) {
	s, _ := newTopTracksTestServer(t)
	req := httptest.NewRequest(http.MethodGet,
		"/top-tracks?artist_id=artist-1&start=notadate&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status %d, want 400", w.Code)
	}
}

func TestGetTopTracksHandler_ServiceError(t *testing.T) {
	s, svc := newTopTracksTestServer(t)
	svc.EXPECT().
		GetTopTracks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet,
		"/top-tracks?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status %d, want 500", w.Code)
	}
}

func TestGetTopTracksHandler_EmptyResult(t *testing.T) {
	s, svc := newTopTracksTestServer(t)
	svc.EXPECT().
		GetTopTracks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]lib.TrackPlay{}, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/top-tracks?artist_id=artist-1&start=2022-01-01T00:00:00Z&end=2022-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", w.Code)
	}
	var resp struct {
		Tracks []lib.TrackPlay `json:"tracks"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Tracks == nil || len(resp.Tracks) != 0 {
		t.Errorf("expected empty array, got %+v", resp.Tracks)
	}
}
