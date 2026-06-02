//go:build integration

package clickhouse_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	goch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/el10savio96/artistInsights/src/infra/clickhouse"
)

var store *clickhouse.ArtistStore

const (
	testArtistID = "test-artist-integration"
	testUserA    = "user-a"
	testUserB    = "user-b"
)

var (
	seedStart = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	seedEnd   = time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)
	seedRows  = []struct {
		artistID, trackID, trackName, userID, artistName string
		ts                                                time.Time
	}{
		{testArtistID, "track-1", "Song One", testUserA, "Test Artist", time.Date(2022, 3, 1, 12, 0, 0, 0, time.UTC)},
		{testArtistID, "track-2", "Song Two", testUserA, "Test Artist", time.Date(2022, 6, 1, 12, 0, 0, 0, time.UTC)},
		{testArtistID, "track-3", "Song Three", testUserB, "Test Artist", time.Date(2022, 9, 1, 12, 0, 0, 0, time.UTC)},
	}
)

func chHost() string {
	if h := os.Getenv("CH_HOST"); h != "" {
		return h
	}
	return "localhost"
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	conn, err := goch.Open(&goch.Options{
		Addr: []string{fmt.Sprintf("%s:9000", chHost())},
		Auth: goch.Auth{
			Database: "artist_insights",
			Username: "admin",
			Password: "pass",
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "clickhouse open: %v\n", err)
		os.Exit(1)
	}
	if err := conn.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "clickhouse ping: %v\n", err)
		os.Exit(1)
	}

	store = &clickhouse.ArtistStore{Conn: conn}

	// seed
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO listens(artist_id, track_id, track_name, user_id, ts, artist_name)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare batch: %v\n", err)
		os.Exit(1)
	}
	for _, r := range seedRows {
		if err := batch.Append(r.artistID, r.trackID, r.trackName, r.userID, r.ts, r.artistName); err != nil {
			fmt.Fprintf(os.Stderr, "append row: %v\n", err)
			os.Exit(1)
		}
	}
	if err := batch.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "send batch: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	// cleanup
	conn.Exec(ctx, "DELETE FROM listens WHERE artist_id = ?", testArtistID)
	conn.Close()
	os.Exit(code)
}

func TestGetSummary_FullRange(t *testing.T) {
	got, err := store.GetSummary(context.Background(), testArtistID, seedStart, seedEnd)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if got.TotalPlays != 3 {
		t.Errorf("TotalPlays = %d, want 3", got.TotalPlays)
	}
	if got.TotalUniqueListeners != 2 {
		t.Errorf("TotalUniqueListeners = %d, want 2", got.TotalUniqueListeners)
	}
}

func TestGetSummary_EmptyRange(t *testing.T) {
	before := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := store.GetSummary(context.Background(), testArtistID, before, before)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if got.TotalPlays != 0 {
		t.Errorf("TotalPlays = %d, want 0", got.TotalPlays)
	}
	if got.TotalUniqueListeners != 0 {
		t.Errorf("TotalUniqueListeners = %d, want 0", got.TotalUniqueListeners)
	}
}

func TestGetSummary_SingleUser(t *testing.T) {
	// Range covering only userA's first play
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 4, 30, 23, 59, 59, 0, time.UTC)

	got, err := store.GetSummary(context.Background(), testArtistID, start, end)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if got.TotalPlays != 1 {
		t.Errorf("TotalPlays = %d, want 1", got.TotalPlays)
	}
	if got.TotalUniqueListeners != 1 {
		t.Errorf("TotalUniqueListeners = %d, want 1", got.TotalUniqueListeners)
	}
}
