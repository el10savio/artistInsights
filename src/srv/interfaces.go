package srv

import (
	"context"
	"time"

	"github.com/el10savio96/artistInsights/src/lib"
)

// ArtistService is the lib contract that srv depends on.
type ArtistService interface {
	GetArtists(ctx context.Context) ([]lib.Artist, error)
	GetSummary(ctx context.Context, artistID string, start, end time.Time) (*lib.Summary, error)
	GetSeries(ctx context.Context, artistID string, start, end time.Time, granularity string) ([]lib.SeriesPoint, error)
	GetHistogram(ctx context.Context, artistID string, start, end time.Time, by string) ([]lib.HistogramBin, error)
	GetTopTracks(ctx context.Context, artistID string, start, end time.Time) ([]lib.TrackPlay, error)
	GetHeatmap(ctx context.Context, artistID string, start, end time.Time) ([]lib.HeatmapCell, error)
	GetWorldMap(ctx context.Context, artistID string, start, end time.Time) ([]lib.CountryListens, error)
}
