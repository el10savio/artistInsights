package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/el10savio96/artistInsights/src/lib"
)

type ArtistStore struct {
	Conn clickhouse.Conn
}

func (s *ArtistStore) GetArtists(ctx context.Context) ([]lib.Artist, error) {
	rows, err := s.Conn.Query(ctx, "SELECT artist_id, artist_name FROM artists FINAL ORDER BY artist_name")
	if err != nil {
		return nil, fmt.Errorf("query artists: %w", err)
	}
	defer rows.Close()

	var artists []lib.Artist
	for rows.Next() {
		var a lib.Artist
		if err := rows.Scan(&a.ArtistID, &a.ArtistName); err != nil {
			return nil, fmt.Errorf("scan artist: %w", err)
		}
		artists = append(artists, a)
	}

	return artists, rows.Err()
}

func (s *ArtistStore) GetSeries(ctx context.Context, artistID string, start, end time.Time, granularity string) ([]lib.SeriesPoint, error) {
	query := fmt.Sprintf(
		`SELECT DATE_TRUNC('%s', ts) AS bucket, count(), COUNT(DISTINCT user_id)
		 FROM listens
		 WHERE artist_id = ? AND ts >= ? AND ts <= ?
		 GROUP BY bucket
		 ORDER BY bucket`,
		granularity, // validated: one of hour|day|week|month
	)
	rows, err := s.Conn.Query(ctx, query, artistID, start, end)
	if err != nil {
		return nil, fmt.Errorf("query series: %w", err)
	}
	defer rows.Close()

	points := []lib.SeriesPoint{}
	for rows.Next() {
		var p lib.SeriesPoint
		if err := rows.Scan(&p.Bucket, &p.TotalPlays, &p.TotalUniqueListeners); err != nil {
			return nil, fmt.Errorf("scan series point: %w", err)
		}
		points = append(points, p)
	}
	return points, rows.Err()
}

var entityCol = map[string]string{
	"listener": "user_id",
	"track":    "track_id",
}

var histogramBinLabels = []string{"1–9", "10–99", "100–999", "1000+"}

func (s *ArtistStore) GetHistogram(ctx context.Context, artistID string, start, end time.Time, by string) ([]lib.HistogramBin, error) {
	query := fmt.Sprintf(
		`WITH per_entity AS (
		     SELECT %s, count() AS play_count
		     FROM listens
		     WHERE artist_id = ? AND ts >= ? AND ts <= ?
		     GROUP BY %s
		 )
		 SELECT
		     multiIf(play_count < 10,   '1–9',
		             play_count < 100,  '10–99',
		             play_count < 1000, '100–999',
		                                '1000+') AS bin,
		     multiIf(play_count < 10,   1,
		             play_count < 100,  2,
		             play_count < 1000, 3,
		                                4)       AS bin_order,
		     count() AS entity_count
		 FROM per_entity
		 GROUP BY bin, bin_order
		 ORDER BY bin_order`,
		entityCol[by], entityCol[by], // validated: listener|track
	)
	rows, err := s.Conn.Query(ctx, query, artistID, start, end)
	if err != nil {
		return nil, fmt.Errorf("query histogram: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]uint64)
	for rows.Next() {
		var bin string
		var binOrder uint8
		var count uint64
		if err := rows.Scan(&bin, &binOrder, &count); err != nil {
			return nil, fmt.Errorf("scan histogram bin: %w", err)
		}
		counts[bin] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	bins := make([]lib.HistogramBin, 0, len(histogramBinLabels))
	for _, label := range histogramBinLabels {
		bins = append(bins, lib.HistogramBin{Bin: label, Count: counts[label]})
	}
	return bins, nil
}

func (s *ArtistStore) GetTopTracks(ctx context.Context, artistID string, start, end time.Time) ([]lib.TrackPlay, error) {
	rows, err := s.Conn.Query(ctx,
		`SELECT track_id, any(track_name) AS track_name, count() AS plays
		 FROM listens
		 WHERE artist_id = ? AND ts >= ? AND ts <= ?
		 GROUP BY track_id
		 ORDER BY plays DESC
		 LIMIT 5`,
		artistID, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("query top tracks: %w", err)
	}
	defer rows.Close()

	tracks := []lib.TrackPlay{}
	for rows.Next() {
		var t lib.TrackPlay
		if err := rows.Scan(&t.TrackID, &t.TrackName, &t.Plays); err != nil {
			return nil, fmt.Errorf("scan track play: %w", err)
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}

func (s *ArtistStore) GetHeatmap(ctx context.Context, artistID string, start, end time.Time) ([]lib.HeatmapCell, error) {
	rows, err := s.Conn.Query(ctx,
		`SELECT toDayOfWeek(ts) AS dow, toHour(ts) AS hour, count() AS plays
		 FROM listens
		 WHERE artist_id = ? AND ts >= ? AND ts <= ?
		 GROUP BY dow, hour
		 ORDER BY dow, hour`,
		artistID, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("query heatmap: %w", err)
	}
	defer rows.Close()

	cells := []lib.HeatmapCell{}
	for rows.Next() {
		var c lib.HeatmapCell
		if err := rows.Scan(&c.Dow, &c.Hour, &c.Plays); err != nil {
			return nil, fmt.Errorf("scan heatmap cell: %w", err)
		}
		cells = append(cells, c)
	}
	return cells, rows.Err()
}

func (s *ArtistStore) GetWorldMap(ctx context.Context, artistID string, start, end time.Time) ([]lib.CountryListens, error) {
	rows, err := s.Conn.Query(ctx,
		`SELECT u.country, count() AS total_listens
		 FROM artist_insights.listens l
		 JOIN artist_insights.users u ON l.user_id = u.user_id
		 WHERE l.artist_id = ? AND l.ts >= ? AND l.ts <= ?
		 GROUP BY u.country
		 ORDER BY total_listens DESC`,
		artistID, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("query world map: %w", err)
	}
	defer rows.Close()

	result := []lib.CountryListens{}
	for rows.Next() {
		var cl lib.CountryListens
		if err := rows.Scan(&cl.Country, &cl.Count); err != nil {
			return nil, fmt.Errorf("scan country listens: %w", err)
		}
		result = append(result, cl)
	}
	return result, rows.Err()
}

func (s *ArtistStore) GetSummary(ctx context.Context, artistID string, start, end time.Time) (*lib.Summary, error) {
	row := s.Conn.QueryRow(ctx,
		`SELECT count(), COUNT(DISTINCT user_id)
		 FROM listens
		 WHERE artist_id = ? AND ts >= ? AND ts <= ?`,
		artistID, start, end,
	)
	var sum lib.Summary
	if err := row.Scan(&sum.TotalPlays, &sum.TotalUniqueListeners); err != nil {
		return nil, fmt.Errorf("query summary: %w", err)
	}
	return &sum, nil
}
