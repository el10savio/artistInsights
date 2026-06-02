package lib

import (
	"time"
)

type Artist struct {
	ArtistID   string `json:"artist_id"`
	ArtistName string `json:"artist_name"`
}

type Summary struct {
	TotalPlays           uint64 `json:"total_plays"`
	TotalUniqueListeners uint64 `json:"total_unique_listeners"`
}

type SeriesPoint struct {
	Bucket               time.Time `json:"bucket"`
	TotalPlays           uint64    `json:"total_plays"`
	TotalUniqueListeners uint64    `json:"total_unique_listeners"`
}

type HistogramBin struct {
	Bin   string `json:"bin"`
	Count uint64 `json:"count"`
}

type TrackPlay struct {
	TrackID   string `json:"track_id"`
	TrackName string `json:"track_name"`
	Plays     uint64 `json:"plays"`
}

type HeatmapCell struct {
	Dow   uint8  `json:"dow"`
	Hour  uint8  `json:"hour"`
	Plays uint64 `json:"plays"`
}

type CountryListens struct {
	Country string `json:"country"`
	Count   uint64 `json:"count"`
}
