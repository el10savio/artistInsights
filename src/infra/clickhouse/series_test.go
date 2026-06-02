//go:build integration

package clickhouse_test

import (
	"context"
	"testing"
	"time"
)

func TestGetSeries_DayGranularity(t *testing.T) {
	// 3 seed rows on 3 distinct days → 3 buckets
	got, err := store.GetSeries(context.Background(), testArtistID, seedStart, seedEnd, "day")
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len(series) = %d, want 3", len(got))
	}
	// Each day has exactly 1 play
	for _, p := range got {
		if p.TotalPlays != 1 {
			t.Errorf("bucket %s: TotalPlays = %d, want 1", p.Bucket, p.TotalPlays)
		}
	}
}

func TestGetSeries_MonthGranularity(t *testing.T) {
	// Rows fall in March, June, September → 3 monthly buckets
	got, err := store.GetSeries(context.Background(), testArtistID, seedStart, seedEnd, "month")
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len(series) = %d, want 3", len(got))
	}

	wantMonths := []time.Month{time.March, time.June, time.September}
	for i, p := range got {
		if p.Bucket.Month() != wantMonths[i] {
			t.Errorf("bucket %d: month = %s, want %s", i, p.Bucket.Month(), wantMonths[i])
		}
		if p.TotalPlays != 1 {
			t.Errorf("bucket %d: TotalPlays = %d, want 1", i, p.TotalPlays)
		}
	}
}

func TestGetSeries_EmptyRange(t *testing.T) {
	before := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := store.GetSeries(context.Background(), testArtistID, before, before, "day")
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty series, got %d points", len(got))
	}
}
