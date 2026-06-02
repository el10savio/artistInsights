import { useState, useEffect, useRef } from 'react'
import './App.css'
import { getSummary, getSeries, getHistogram, getTopTracks, getHeatmap, getWorldMap } from './api/client'
import type { Summary, SeriesPoint, HistogramBin, TrackPlay, HeatmapCell, CountryListens, Granularity } from './types/api'
import ArtistSelector from './components/ArtistSelector'
import DateRangePicker from './components/DateRangePicker'
import GranularitySelector from './components/GranularitySelector'
import SummaryTiles from './components/SummaryTiles'
import TimeSeriesChart from './components/TimeSeriesChart'
import ListenerHistogram from './components/HistogramCharts'
import TopTracksChart from './components/TopTracksChart'
import HeatmapChart from './components/HeatmapChart'
import WorldMapChart from './components/WorldMapChart'
import HelpTooltip from './components/HelpTooltip'

const DATASET_START = new Date('2008-01-27T00:00:00Z')
const DATASET_END = new Date('2009-05-04T23:59:59Z')

export default function App() {
  const [artistId, setArtistId] = useState<string | null>(null)
  const [start, setStart] = useState<Date>(DATASET_START)
  const [end, setEnd] = useState<Date>(DATASET_END)
  const [granularity, setGranularity] = useState<Granularity>('day')

  const [summary, setSummary] = useState<Summary | null>(null)
  const [summaryLoading, setSummaryLoading] = useState(false)

  const [series, setSeries] = useState<SeriesPoint[]>([])
  const [seriesLoading, setSeriesLoading] = useState(false)

  const [listenerHist, setListenerHist] = useState<HistogramBin[]>([])
  const [listenerLoading, setListenerLoading] = useState(false)

  const [topTracks, setTopTracks] = useState<TrackPlay[]>([])
  const [topTracksLoading, setTopTracksLoading] = useState(false)

  const [heatmap, setHeatmap] = useState<HeatmapCell[]>([])
  const [heatmapLoading, setHeatmapLoading] = useState(false)

  const [worldMap, setWorldMap] = useState<CountryListens[]>([])
  const [worldMapLoading, setWorldMapLoading] = useState(false)

  const seriesAbort = useRef<AbortController | null>(null)
  const histAbort = useRef<AbortController | null>(null)

  useEffect(() => {
    if (!artistId) {
      setSummary(null)
      setSeries([])
      return
    }

    seriesAbort.current?.abort()
    seriesAbort.current = new AbortController()

    setSummaryLoading(true)
    setSeriesLoading(true)

    Promise.all([
      getSummary(artistId, start, end),
      getSeries(artistId, start, end, granularity),
    ])
      .then(([s, pts]) => {
        setSummary(s)
        setSeries(pts)
      })
      .catch(() => {})
      .finally(() => {
        setSummaryLoading(false)
        setSeriesLoading(false)
      })
  }, [artistId, start, end, granularity])

  useEffect(() => {
    if (!artistId) {
      setListenerHist([])
      setTopTracks([])
      setHeatmap([])
      return
    }

    histAbort.current?.abort()
    histAbort.current = new AbortController()

    setListenerLoading(true)
    setTopTracksLoading(true)
    setHeatmapLoading(true)

    Promise.all([
      getHistogram(artistId, start, end, 'listener'),
      getTopTracks(artistId, start, end),
      getHeatmap(artistId, start, end),
    ])
      .then(([lh, tt, hm]) => {
        setListenerHist(lh)
        setTopTracks(tt)
        setHeatmap(hm)
      })
      .catch(() => {})
      .finally(() => {
        setListenerLoading(false)
        setTopTracksLoading(false)
        setHeatmapLoading(false)
      })
  }, [artistId, start, end])

  useEffect(() => {
    if (!artistId) {
      setWorldMap([])
      return
    }
    setWorldMapLoading(true)
    getWorldMap(artistId, start, end)
      .then(setWorldMap)
      .catch(() => {})
      .finally(() => setWorldMapLoading(false))
  }, [artistId, start, end])

  const hasArtist = artistId !== null

  return (
    <div className="app">
      <div className="app-header">
        <h1>Artist Insights</h1>
        <p>Last.fm listening data · Jan 2008 – May 2009</p>
      </div>

      <div className="controls">
        <div className="control-group" style={{ flex: '2 1 260px' }}>
          <label>Artist</label>
          <ArtistSelector value={artistId} onChange={setArtistId} />
        </div>
        <div className="control-group">
          <label>Date range</label>
          <DateRangePicker
            start={start}
            end={end}
            onStartChange={setStart}
            onEndChange={setEnd}
          />
        </div>
        <div className="control-group" style={{ flex: '0 0 auto' }}>
          <label>Granularity</label>
          <GranularitySelector value={granularity} onChange={setGranularity} />
        </div>
      </div>

      {!hasArtist ? (
        <div className="prompt">
          <span className="prompt-icon">↑</span>
          Select an artist to explore their listening data
        </div>
      ) : (
        <>
          <SummaryTiles data={summary} loading={summaryLoading} />

          <div className="histograms" style={{ marginBottom: 24 }}>
            <div className="chart-card" style={{ marginBottom: 0 }}>
              <h2>Plays Over Time<HelpTooltip text="Play count per time bucket. Use the granularity selector to change the bucket size." /></h2>
              <TimeSeriesChart
                data={series}
                loading={seriesLoading}
                granularity={granularity}
                dataKey="total_plays"
                color="#37bbe4"
                label="Plays"
              />
            </div>
            <div className="chart-card" style={{ marginBottom: 0 }}>
              <h2>Unique Listeners Over Time<HelpTooltip text="Distinct listener count per time bucket. A listener is counted once per bucket regardless of how many tracks they played." /></h2>
              <TimeSeriesChart
                data={series}
                loading={seriesLoading}
                granularity={granularity}
                dataKey="total_unique_listeners"
                color="#a8d5e8"
                label="Unique Listeners"
              />
            </div>
          </div>

          <div className="histograms" style={{ marginBottom: 24 }}>
            <ListenerHistogram data={listenerHist} loading={listenerLoading} />
            <TopTracksChart data={topTracks} loading={topTracksLoading} />
          </div>

          <HeatmapChart data={heatmap} loading={heatmapLoading} />

          <div className="chart-card" style={{ marginTop: 24 }}>
            <h2>
              Total Listens by Country
              <HelpTooltip text="Total listens sourced from each country across all artists for the selected date range." />
            </h2>
            <WorldMapChart data={worldMap} loading={worldMapLoading} />
          </div>
        </>
      )}
    </div>
  )
}
