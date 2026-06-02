import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { SeriesPoint, Granularity } from '../types/api'

interface Props {
  data: SeriesPoint[]
  loading: boolean
  granularity: Granularity
  dataKey: 'total_plays' | 'total_unique_listeners'
  color: string
  label: string
}

function formatBucket(bucket: string, granularity: Granularity): string {
  const d = new Date(bucket)
  if (granularity === 'hour') {
    return d.toLocaleString('en-GB', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'UTC' })
  }
  if (granularity === 'month') {
    return d.toLocaleString('en-GB', { month: 'short', year: '2-digit', timeZone: 'UTC' })
  }
  if (granularity === 'week') {
    return d.toLocaleString('en-GB', { month: 'short', day: 'numeric', year: '2-digit', timeZone: 'UTC' })
  }
  return d.toLocaleString('en-GB', { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

export default function TimeSeriesChart({ data, loading, granularity, dataKey, color, label }: Props) {
  if (loading) {
    return (
      <div className="empty-state">
        <span className="spinner" />
        Loading…
      </div>
    )
  }

  if (data.length === 0) {
    return <div className="empty-state">No data for this range</div>
  }

  const chartData = data.map((p) => ({
    ...p,
    label: formatBucket(p.bucket, granularity),
  }))

  return (
    <ResponsiveContainer width="100%" height={220}>
      <LineChart data={chartData} margin={{ top: 4, right: 12, bottom: 0, left: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="rgba(225,224,221,0.2)" />
        <XAxis
          dataKey="label"
          tick={{ fill: '#6b6760', fontSize: 11 }}
          tickLine={false}
          axisLine={{ stroke: 'rgba(225,224,221,0.2)' }}
          interval="equidistantPreserveStart"
        />
        <YAxis
          tick={{ fill: '#6b6760', fontSize: 11 }}
          tickLine={false}
          axisLine={false}
          width={48}
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
        />
        <Line
          type="monotone"
          dataKey={dataKey}
          name={label}
          stroke={color}
          strokeWidth={2}
          dot={false}
          activeDot={{ r: 4 }}
        />
      </LineChart>
    </ResponsiveContainer>
  )
}
