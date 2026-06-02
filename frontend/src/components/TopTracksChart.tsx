import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { TrackPlay } from '../types/api'
import HelpTooltip from './HelpTooltip'

interface Props {
  data: TrackPlay[]
  loading: boolean
}

export default function TopTracksChart({ data, loading }: Props) {
  const chartData = data.map((t) => ({
    label: t.track_name,
    plays: t.plays,
    fullId: t.track_id,
  }))

  return (
    <div className="chart-card" style={{ marginBottom: 0 }}>
      <h2>
        Top Tracks
        <HelpTooltip text="The 5 most-played tracks by this artist in the selected date range, ranked by total play count." />
      </h2>
      {loading ? (
        <div className="empty-state">
          <span className="spinner" />
          Loading…
        </div>
      ) : data.length === 0 ? (
        <div className="empty-state">No data</div>
      ) : (
        <ResponsiveContainer width="100%" height={220}>
          <BarChart layout="vertical" data={chartData} margin={{ top: 4, right: 16, bottom: 0, left: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="rgba(53,52,47,0.15)" horizontal={false} />
            <XAxis
              type="number"
              tick={{ fill: '#6b6760', fontSize: 10 }}
              tickLine={false}
              axisLine={{ stroke: 'rgba(53,52,47,0.15)' }}
            />
            <YAxis
              type="category"
              dataKey="label"
              tick={{ fill: '#6b6760', fontSize: 10 }}
              tickLine={false}
              axisLine={false}
              width={160}
            />
            <Tooltip
              contentStyle={{
                background: '#2a2924',
                border: '1px solid rgba(225,224,221,0.2)',
                borderRadius: 6,
                fontSize: 12,
              }}
              labelStyle={{ color: '#f1f2f0', marginBottom: 4 }}
              itemStyle={{ color: '#f1f2f0' }}
              labelFormatter={(_label, payload) =>
                payload?.[0]?.payload?.fullId ?? _label
              }
            />
            <Bar dataKey="plays" name="Plays" fill="#37bbe4" radius={[0, 3, 3, 0]} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
