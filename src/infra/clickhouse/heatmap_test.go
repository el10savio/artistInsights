//go:build integration

package clickhouse_test

import (
	"context"
	"testing"
	"time"
)

// Seed rows and their expected (dow, hour):
//   2022-03-01 12:00 UTC → Tuesday  (dow=2), hour=12
//   2022-06-01 12:00 UTC → Wednesday(dow=3), hour=12
//   2022-09-01 12:00 UTC → Thursday (dow=4), hour=12

func TestGetHeatmap_ReturnsSeedCells(t *testing.T) {
	cells, err := store.GetHeatmap(context.Background(), testArtistID, seedStart, seedEnd)
	if err != nil {
		t.Fatalf("GetHeatmap: %v", err)
	}
	if len(cells) != 3 {
		t.Fatalf("expected 3 cells, got %d: %+v", len(cells), cells)
	}

	type key struct{ dow, hour uint8 }
	got := make(map[key]uint64)
	for _, c := range cells {
		got[key{c.Dow, c.Hour}] = c.Plays
	}

	expected := map[key]uint64{
		{2, 12}: 1,
		{3, 12}: 1,
		{4, 12}: 1,
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("cell dow=%d hour=%d: plays=%d, want %d", k.dow, k.hour, got[k], want)
		}
	}
}

func TestGetHeatmap_EmptyRange(t *testing.T) {
	before := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	cells, err := store.GetHeatmap(context.Background(), testArtistID, before, before)
	if err != nil {
		t.Fatalf("GetHeatmap: %v", err)
	}
	if len(cells) != 0 {
		t.Errorf("expected empty, got %d cells", len(cells))
	}
}
