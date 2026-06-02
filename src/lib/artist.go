package lib

import (
	"context"
	"time"
)

// ArtistRepository is the infra contract that lib depends on.
type ArtistRepository interface {
	GetArtists(ctx context.Context) ([]Artist, error)
	GetSummary(ctx context.Context, artistID string, start, end time.Time) (*Summary, error)
	GetSeries(ctx context.Context, artistID string, start, end time.Time, granularity string) ([]SeriesPoint, error)
	GetHistogram(ctx context.Context, artistID string, start, end time.Time, by string) ([]HistogramBin, error)
	GetTopTracks(ctx context.Context, artistID string, start, end time.Time) ([]TrackPlay, error)
	GetHeatmap(ctx context.Context, artistID string, start, end time.Time) ([]HeatmapCell, error)
	GetWorldMap(ctx context.Context, artistID string, start, end time.Time) ([]CountryListens, error)
}

type ArtistService struct {
	repo ArtistRepository
}

func NewArtistService(repo ArtistRepository) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) GetArtists(ctx context.Context) ([]Artist, error) {
	return s.repo.GetArtists(ctx)
}

func (s *ArtistService) GetSummary(ctx context.Context, artistID string, start, end time.Time) (*Summary, error) {
	return s.repo.GetSummary(ctx, artistID, start, end)
}

func (s *ArtistService) GetSeries(ctx context.Context, artistID string, start, end time.Time, granularity string) ([]SeriesPoint, error) {
	return s.repo.GetSeries(ctx, artistID, start, end, granularity)
}

func (s *ArtistService) GetHistogram(ctx context.Context, artistID string, start, end time.Time, by string) ([]HistogramBin, error) {
	return s.repo.GetHistogram(ctx, artistID, start, end, by)
}

func (s *ArtistService) GetTopTracks(ctx context.Context, artistID string, start, end time.Time) ([]TrackPlay, error) {
	return s.repo.GetTopTracks(ctx, artistID, start, end)
}

func (s *ArtistService) GetHeatmap(ctx context.Context, artistID string, start, end time.Time) ([]HeatmapCell, error) {
	return s.repo.GetHeatmap(ctx, artistID, start, end)
}

func (s *ArtistService) GetWorldMap(ctx context.Context, artistID string, start, end time.Time) ([]CountryListens, error) {
	return s.repo.GetWorldMap(ctx, artistID, start, end)
}
