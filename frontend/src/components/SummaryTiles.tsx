import type { Summary } from '../types/api'
import HelpTooltip from './HelpTooltip'

interface Props {
  data: Summary | null
  loading: boolean
}

function fmt(n: number): string {
  return n.toLocaleString()
}

export default function SummaryTiles({ data, loading }: Props) {
  const plays = loading ? '—' : data ? fmt(data.total_plays) : '—'
  const listeners = loading ? '—' : data ? fmt(data.total_unique_listeners) : '—'

  return (
    <div className="summary-tiles">
      <div className="tile">
        <div className="tile-label">Total Plays<HelpTooltip text="Total number of plays for this artist in the selected date range." /></div>
        <div className="tile-value" style={loading ? { opacity: 0.4 } : {}}>
          {plays}
        </div>
      </div>
      <div className="tile">
        <div className="tile-label">Unique Listeners<HelpTooltip text="Number of distinct users who played at least one track in the selected date range." /></div>
        <div className="tile-value" style={loading ? { opacity: 0.4 } : {}}>
          {listeners}
        </div>
      </div>
    </div>
  )
}
