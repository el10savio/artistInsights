import type { Artist, Summary, SeriesPoint, HistogramBin, TrackPlay, HeatmapCell, CountryListens, Granularity, HistogramBy } from '../types/api'

function buildUrl(path: string, params: Record<string, string>): string {
  const qs = new URLSearchParams(params).toString()
  return `/api${path}?${qs}`
}

async function get<T>(url: string): Promise<T> {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  return res.json() as Promise<T>
}

export async function getArtists(): Promise<Artist[]> {
  const data = await get<{ artists: Artist[] }>('/api/artists')
  return data.artists ?? []
}

export async function getSummary(
  artistId: string,
  start: Date,
  end: Date,
): Promise<Summary> {
  const url = buildUrl('/summary', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
  })
  return get<Summary>(url)
}

export async function getSeries(
  artistId: string,
  start: Date,
  end: Date,
  granularity: Granularity,
): Promise<SeriesPoint[]> {
  const url = buildUrl('/series', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
    granularity,
  })
  const data = await get<{ series: SeriesPoint[] }>(url)
  return data.series ?? []
}

export async function getHeatmap(
  artistId: string,
  start: Date,
  end: Date,
): Promise<HeatmapCell[]> {
  const url = buildUrl('/heatmap', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
  })
  const data = await get<{ heatmap: HeatmapCell[] }>(url)
  return data.heatmap ?? []
}

export async function getTopTracks(
  artistId: string,
  start: Date,
  end: Date,
): Promise<TrackPlay[]> {
  const url = buildUrl('/top-tracks', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
  })
  const data = await get<{ tracks: TrackPlay[] }>(url)
  return data.tracks ?? []
}

export async function getHistogram(
  artistId: string,
  start: Date,
  end: Date,
  by: HistogramBy,
): Promise<HistogramBin[]> {
  const url = buildUrl('/histogram', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
    by,
  })
  const data = await get<{ histogram: HistogramBin[] }>(url)
  return data.histogram ?? []
}

export async function getWorldMap(artistId: string, start: Date, end: Date): Promise<CountryListens[]> {
  const url = buildUrl('/world-map', {
    artist_id: artistId,
    start: start.toISOString(),
    end: end.toISOString(),
  })
  const data = await get<{ world_map: CountryListens[] }>(url)
  return data.world_map ?? []
}
