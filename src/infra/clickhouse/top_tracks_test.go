//go:build integration

package clickhouse_test

import (
	"context"
	"testing"
)

func TestGetTopTracks_ReturnsOrdered(t *testing.T) {
	// Seed: user-a played track-1 and track-2, user-b played track-3 → track-1 and track-2 each have 1 play
	tracks, err := store.GetTopTracks(context.Background(), testArtistID, seedStart, seedEnd)
	if err != nil {
		t.Fatalf("GetTopTracks: %v", err)
	}
	if len(tracks) == 0 {
		t.Fatal("expected tracks, got none")
	}
	for i := 1; i < len(tracks); i++ {
		if tracks[i].Plays > tracks[i-1].Plays {
			t.Errorf("tracks not sorted: [%d]=%d > [%d]=%d",
				i, tracks[i].Plays, i-1, tracks[i-1].Plays)
		}
	}
	if len(tracks) > 5 {
		t.Errorf("expected at most 5 tracks, got %d", len(tracks))
	}
}

func TestGetTopTracks_EmptyRange(t *testing.T) {
	before := seedStart.Add(-1)
	tracks, err := store.GetTopTracks(context.Background(), testArtistID, before, before)
	if err != nil {
		t.Fatalf("GetTopTracks: %v", err)
	}
	if len(tracks) != 0 {
		t.Errorf("expected empty, got %d tracks", len(tracks))
	}
}
