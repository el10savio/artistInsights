package lib_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/el10savio/artistInsights/src/lib"
	"github.com/el10savio/artistInsights/src/lib/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetSummary_PropagatesResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	want := &lib.Summary{TotalPlays: 100, TotalUniqueListeners: 42}

	repo.EXPECT().
		GetSummary(gomock.Any(), "artist-1", start, end).
		Return(want, nil)

	svc := lib.NewArtistService(repo)
	got, err := svc.GetSummary(context.Background(), "artist-1", start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != *want {
		t.Errorf("got %+v, want %+v", *got, *want)
	}
}

func TestGetSummary_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	repoErr := errors.New("db down")

	repo.EXPECT().
		GetSummary(gomock.Any(), "artist-1", start, end).
		Return(nil, repoErr)

	svc := lib.NewArtistService(repo)
	_, err := svc.GetSummary(context.Background(), "artist-1", start, end)
	if !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}

func TestGetSeries_PropagatesResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	want := []lib.SeriesPoint{
		{Bucket: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC), TotalPlays: 5, TotalUniqueListeners: 2},
		{Bucket: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC), TotalPlays: 3, TotalUniqueListeners: 1},
	}

	repo.EXPECT().
		GetSeries(gomock.Any(), "artist-1", start, end, "month").
		Return(want, nil)

	svc := lib.NewArtistService(repo)
	got, err := svc.GetSeries(context.Background(), "artist-1", start, end, "month")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d points, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("point %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestGetHistogram_PropagatesResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	want := []lib.HistogramBin{
		{Bin: "1–9", Count: 80},
		{Bin: "10–99", Count: 20},
	}

	repo.EXPECT().
		GetHistogram(gomock.Any(), "artist-1", start, end, "listener").
		Return(want, nil)

	svc := lib.NewArtistService(repo)
	got, err := svc.GetHistogram(context.Background(), "artist-1", start, end, "listener")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d bins, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("bin %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestGetTopTracks_PropagatesResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	want := []lib.TrackPlay{
		{TrackID: "track-1", Plays: 50},
		{TrackID: "track-2", Plays: 30},
	}

	repo.EXPECT().
		GetTopTracks(gomock.Any(), "artist-1", start, end).
		Return(want, nil)

	svc := lib.NewArtistService(repo)
	got, err := svc.GetTopTracks(context.Background(), "artist-1", start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetTopTracks_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	repoErr := errors.New("db down")

	repo.EXPECT().
		GetTopTracks(gomock.Any(), "artist-1", start, end).
		Return(nil, repoErr)

	svc := lib.NewArtistService(repo)
	_, err := svc.GetTopTracks(context.Background(), "artist-1", start, end)
	if !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}

func TestGetHeatmap_PropagatesResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	want := []lib.HeatmapCell{
		{Dow: 1, Hour: 8, Plays: 45},
		{Dow: 5, Hour: 22, Plays: 91},
	}

	repo.EXPECT().
		GetHeatmap(gomock.Any(), "artist-1", start, end).
		Return(want, nil)

	svc := lib.NewArtistService(repo)
	got, err := svc.GetHeatmap(context.Background(), "artist-1", start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != want[0] {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetHeatmap_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	repoErr := errors.New("db down")

	repo.EXPECT().
		GetHeatmap(gomock.Any(), "artist-1", start, end).
		Return(nil, repoErr)

	svc := lib.NewArtistService(repo)
	_, err := svc.GetHeatmap(context.Background(), "artist-1", start, end)
	if !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}

func TestGetHistogram_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	repoErr := errors.New("db down")

	repo.EXPECT().
		GetHistogram(gomock.Any(), "artist-1", start, end, "track").
		Return(nil, repoErr)

	svc := lib.NewArtistService(repo)
	_, err := svc.GetHistogram(context.Background(), "artist-1", start, end, "track")
	if !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}

func TestGetSeries_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockArtistRepository(ctrl)

	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	repoErr := errors.New("db down")

	repo.EXPECT().
		GetSeries(gomock.Any(), "artist-1", start, end, "day").
		Return(nil, repoErr)

	svc := lib.NewArtistService(repo)
	_, err := svc.GetSeries(context.Background(), "artist-1", start, end, "day")
	if !errors.Is(err, repoErr) {
		t.Errorf("got %v, want %v", err, repoErr)
	}
}
