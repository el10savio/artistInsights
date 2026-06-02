import { useState } from 'react'
import type { HeatmapCell } from '../types/api'
import HelpTooltip from './HelpTooltip'

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
const HOURS = Array.from({ length: 24 }, (_, i) => i)

interface Props {
  data: HeatmapCell[]
  loading: boolean
}

interface TooltipState {
  dow: number
  hour: number
  plays: number
  x: number
  y: number
}

export default function HeatmapChart({ data, loading }: Props) {
  const [tooltip, setTooltip] = useState<TooltipState | null>(null)

  // Build lookup: dow → hour → plays
  const grid = new Map<number, Map<number, number>>()
  let maxPlays = 0
  for (const cell of data) {
    if (!grid.has(cell.dow)) grid.set(cell.dow, new Map())
    grid.get(cell.dow)!.set(cell.hour, cell.plays)
    if (cell.plays > maxPlays) maxPlays = cell.plays
  }

  const opacity = (plays: number) =>
    maxPlays === 0 ? 0.05 : Math.max(0.05, plays / maxPlays)

  return (
    <div className="chart-card">
      <h2>
        Listening Heatmap
        <HelpTooltip text="Play count by day of week and hour of day. Shows when this artist's audience is most active." />
      </h2>
      {loading ? (
        <div className="empty-state">
          <span className="spinner" />
          Loading…
        </div>
      ) : data.length === 0 ? (
        <div className="empty-state">No data</div>
      ) : (
        <div style={{ position: 'relative' }}>
          <div className="heatmap-grid">
            {/* hour labels */}
            <div className="heatmap-row-label" />
            {HOURS.map((h) => (
              <div key={h} className="heatmap-hour-label">
                {h % 3 === 0 ? `${h}h` : ''}
              </div>
            ))}

            {/* rows: dow 1–7 */}
            {DAYS.map((day, di) => {
              const dow = di + 1
              return (
                <>
                  <div key={`label-${dow}`} className="heatmap-row-label">{day}</div>
                  {HOURS.map((hour) => {
                    const plays = grid.get(dow)?.get(hour) ?? 0
                    return (
                      <div
                        key={`${dow}-${hour}`}
                        className="heatmap-cell"
                        style={{ background: `rgba(55,187,228,${opacity(plays)})` }}
                        onMouseEnter={(e) => {
                          const rect = (e.target as HTMLElement).getBoundingClientRect()
                          const parent = (e.target as HTMLElement).closest('.chart-card')!.getBoundingClientRect()
                          setTooltip({
                            dow, hour, plays,
                            x: rect.left - parent.left + rect.width / 2,
                            y: rect.top - parent.top,
                          })
                        }}
                        onMouseLeave={() => setTooltip(null)}
                      />
                    )
                  })}
                </>
              )
            })}
          </div>

          {tooltip && (
            <div
              className="heatmap-tooltip"
              style={{ left: tooltip.x, top: tooltip.y - 8 }}
            >
              {DAYS[tooltip.dow - 1]} {String(tooltip.hour).padStart(2, '0')}:00 · {tooltip.plays.toLocaleString()} plays
            </div>
          )}
        </div>
      )}
    </div>
  )
}
