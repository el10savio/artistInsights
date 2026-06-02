export interface Artist {
  artist_id: string
  artist_name: string
}

export interface Summary {
  total_plays: number
  total_unique_listeners: number
}

export interface SeriesPoint {
  bucket: string
  total_plays: number
  total_unique_listeners: number
}

export interface HistogramBin {
  bin: string
  count: number
}

export interface TrackPlay {
  track_id: string
  track_name: string
  plays: number
}

export interface HeatmapCell {
  dow: number   // 1=Mon … 7=Sun
  hour: number  // 0–23
  plays: number
}

export interface CountryListens {
  country: string
  count: number
}

export type Granularity = 'hour' | 'day' | 'week' | 'month'
export type HistogramBy = 'listener' | 'track'
