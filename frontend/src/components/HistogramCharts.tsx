import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { HistogramBin } from '../types/api'
import HelpTooltip from './HelpTooltip'

interface Props {
  data: HistogramBin[]
  loading: boolean
}

export default function ListenerHistogram({ data, loading }: Props) {
  const chartData = data
    .filter((b) => b.count > 0)
    .map((b) => ({
      label: b.bin,
      height: b.count,
    }))

  return (
    <div className="chart-card" style={{ marginBottom: 0 }}>
      <h2>
        Plays per Listener
        <HelpTooltip text="How many plays each listener made. Shows whether the audience is mostly casual (1–9 plays) or contains superfans (100+)." />
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
              width={80}
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
            <Bar dataKey="height" name="Listeners" fill="#37bbe4" radius={[0, 3, 3, 0]} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
