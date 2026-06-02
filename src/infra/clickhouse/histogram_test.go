//go:build integration

package clickhouse_test

import (
	"context"
	"testing"
	"time"
)

func TestGetHistogram_ByListener(t *testing.T) {
	// user-a=2 plays, user-b=1 play → both in "1–9"
	bins, err := store.GetHistogram(context.Background(), testArtistID, seedStart, seedEnd, "listener")
	if err != nil {
		t.Fatalf("GetHistogram: %v", err)
	}
	if len(bins) != 4 {
		t.Fatalf("expected 4 bins, got %d", len(bins))
	}
	expected := map[string]uint64{"1–9": 2, "10–99": 0, "100–999": 0, "1000+": 0}
	for _, b := range bins {
		if b.Count != expected[b.Bin] {
			t.Errorf("bin %q: count=%d, want %d", b.Bin, b.Count, expected[b.Bin])
		}
	}
}

func TestGetHistogram_ByTrack(t *testing.T) {
	// 3 tracks each with 1 play → all in "1–9"
	bins, err := store.GetHistogram(context.Background(), testArtistID, seedStart, seedEnd, "track")
	if err != nil {
		t.Fatalf("GetHistogram: %v", err)
	}
	if len(bins) != 4 {
		t.Fatalf("expected 4 bins, got %d", len(bins))
	}
	expected := map[string]uint64{"1–9": 3, "10–99": 0, "100–999": 0, "1000+": 0}
	for _, b := range bins {
		if b.Count != expected[b.Bin] {
			t.Errorf("bin %q: count=%d, want %d", b.Bin, b.Count, expected[b.Bin])
		}
	}
}

func TestGetHistogram_EmptyRange(t *testing.T) {
	before := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	bins, err := store.GetHistogram(context.Background(), testArtistID, before, before, "listener")
	if err != nil {
		t.Fatalf("GetHistogram: %v", err)
	}
	if len(bins) != 4 {
		t.Fatalf("expected 4 bins, got %d", len(bins))
	}
	for _, b := range bins {
		if b.Count != 0 {
			t.Errorf("bin %q: count=%d, want 0", b.Bin, b.Count)
		}
	}
}
