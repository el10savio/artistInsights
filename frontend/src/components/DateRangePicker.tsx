const MIN_DATE = '2008-01-27'
const MAX_DATE = '2009-05-04'

function toInputVal(d: Date): string {
  return d.toISOString().slice(0, 10)
}

interface Props {
  start: Date
  end: Date
  onStartChange: (d: Date) => void
  onEndChange: (d: Date) => void
}

export default function DateRangePicker({ start, end, onStartChange, onEndChange }: Props) {
  const handleStart = (e: React.ChangeEvent<HTMLInputElement>) => {
    const d = new Date(e.target.value + 'T00:00:00Z')
    if (!isNaN(d.getTime()) && d <= end) onStartChange(d)
  }

  const handleEnd = (e: React.ChangeEvent<HTMLInputElement>) => {
    const d = new Date(e.target.value + 'T23:59:59Z')
    if (!isNaN(d.getTime()) && d >= start) onEndChange(d)
  }

  return (
    <div style={{ display: 'flex', gap: 8, alignItems: 'center', flexWrap: 'wrap' }}>
      <input
        type="date"
        min={MIN_DATE}
        max={toInputVal(end)}
        value={toInputVal(start)}
        onChange={handleStart}
        style={inputStyle}
      />
      <span style={{ color: 'var(--muted)', fontSize: 12 }}>→</span>
      <input
        type="date"
        min={toInputVal(start)}
        max={MAX_DATE}
        value={toInputVal(end)}
        onChange={handleEnd}
        style={inputStyle}
      />
    </div>
  )
}

const inputStyle: React.CSSProperties = {
  padding: '8px 10px',
  borderRadius: 6,
  border: '1px solid rgba(225, 224, 221, 0.3)',
  background: 'rgba(241, 242, 240, 0.08)',
  color: 'var(--surface)',
  fontSize: 13,
  colorScheme: 'dark',
  outline: 'none',
}
