import { useState, useEffect, useMemo } from 'react'
import { geoMercator, geoPath } from 'd3-geo'
import { scaleLinear } from 'd3-scale'
import { feature } from 'topojson-client'
import type { Topology, GeometryCollection } from 'topojson-specification'
import type { Feature, Geometry } from 'geojson'
import type { CountryListens } from '../types/api'

const GEO_URL = 'https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json'

const WIDTH = 800
const HEIGHT = 420

// Maps Last.fm country names → world-atlas TopoJSON names where they differ
const COUNTRY_NAME_MAP: Record<string, string> = {
  'united states': 'united states of america',
  'russian federation': 'russia',
  'czech republic': 'czechia',
  'bosnia and herzegovina': 'bosnia and herz.',
  "cote d'ivoire": "côte d'ivoire",
  'congo, the democratic republic of the': 'dem. rep. congo',
  "korea, democratic people's republic of": 'north korea',
}

interface Props {
  data: CountryListens[]
  loading: boolean
}

interface TooltipState {
  country: string
  count: number
  x: number
  y: number
}

export default function WorldMapChart({ data, loading }: Props) {
  const [geographies, setGeographies] = useState<Feature<Geometry, { name: string }>[]>([])
  const [tooltip, setTooltip] = useState<TooltipState | null>(null)

  useEffect(() => {
    fetch(GEO_URL)
      .then(r => r.json())
      .then((topo: Topology) => {
        const countries = feature(
          topo,
          topo.objects['countries'] as GeometryCollection<{ name: string }>,
        )
        setGeographies(countries.features)
      })
      .catch(() => {})
  }, [])

  // Build lookup keyed by TopoJSON name (lowercased)
  const countryMap = useMemo(() => {
    const m = new Map<string, number>()
    for (const d of data) {
      const key = d.country.toLowerCase()
      const mapped = COUNTRY_NAME_MAP[key] ?? key
      m.set(mapped, d.count)
    }
    return m
  }, [data])

  const maxCount = useMemo(() => Math.max(...data.map(d => d.count), 1), [data])

  const colorScale = useMemo(
    () => scaleLinear<string>().domain([0, maxCount]).range(['#c8eaf5', '#37bbe4']),
    [maxCount],
  )

  const projection = useMemo(
    () => geoMercator().scale(120).center([0, 20]).translate([WIDTH / 2, HEIGHT / 2]),
    [],
  )

  const pathGen = useMemo(() => geoPath(projection), [projection])

  if (loading) {
    return <div className="chart-placeholder">Loading world map…</div>
  }

  return (
    <div style={{ position: 'relative' }}>
      <svg
        viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
        style={{ width: '100%', height: 'auto', display: 'block' }}
      >
        {geographies.map((geo, i) => {
          const name = geo.properties?.name ?? ''
          const count = countryMap.get(name.toLowerCase()) ?? 0
          const fill = count > 0 ? colorScale(count) : '#d4d3cf'
          const d = pathGen(geo) ?? ''
          return (
            <path
              key={i}
              d={d}
              fill={fill}
              stroke="#35342f"
              strokeWidth={0.4}
              style={{ cursor: count > 0 ? 'pointer' : 'default' }}
              onMouseEnter={e => setTooltip({ country: name, count, x: e.clientX, y: e.clientY })}
              onMouseMove={e =>
                setTooltip(prev => (prev ? { ...prev, x: e.clientX, y: e.clientY } : null))
              }
              onMouseLeave={() => setTooltip(null)}
            />
          )
        })}
      </svg>

      {tooltip && (
        <div
          style={{
            position: 'fixed',
            left: tooltip.x + 12,
            top: tooltip.y - 8,
            background: '#35342f',
            color: '#f1f2f0',
            padding: '4px 8px',
            borderRadius: 4,
            fontSize: 12,
            pointerEvents: 'none',
            zIndex: 100,
            whiteSpace: 'nowrap',
          }}
        >
          <strong>{tooltip.country}</strong>
          {': '}
          {tooltip.count > 0 ? tooltip.count.toLocaleString() + ' listens' : 'no data'}
        </div>
      )}
    </div>
  )
}
