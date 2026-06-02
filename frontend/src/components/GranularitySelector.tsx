import type { Granularity } from '../types/api'

const OPTIONS: { value: Granularity; label: string }[] = [
  { value: 'hour', label: 'Hour' },
  { value: 'day', label: 'Day' },
  { value: 'week', label: 'Week' },
  { value: 'month', label: 'Month' },
]

interface Props {
  value: Granularity
  onChange: (g: Granularity) => void
}

export default function GranularitySelector({ value, onChange }: Props) {
  return (
    <div style={wrapStyle}>
      {OPTIONS.map((opt) => (
        <button
          key={opt.value}
          onClick={() => onChange(opt.value)}
          style={{
            ...btnStyle,
            background: value === opt.value ? 'var(--accent)' : 'rgba(241,242,240,0.08)',
            color: value === opt.value ? '#fff' : 'var(--muted)',
            borderColor: value === opt.value ? 'var(--accent)' : 'rgba(225,224,221,0.2)',
          }}
        >
          {opt.label}
        </button>
      ))}
    </div>
  )
}

const wrapStyle: React.CSSProperties = {
  display: 'flex',
  gap: 4,
}

const btnStyle: React.CSSProperties = {
  padding: '7px 12px',
  borderRadius: 6,
  border: '1px solid',
  fontSize: 12,
  fontWeight: 500,
  transition: 'background 0.15s, color 0.15s',
}
