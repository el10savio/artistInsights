import { useState, useEffect, useRef, useCallback } from 'react'
import { getArtists } from '../api/client'
import type { Artist } from '../types/api'

interface Props {
  value: string | null
  onChange: (id: string | null) => void
}

export default function ArtistSelector({ value, onChange }: Props) {
  const [allArtists, setAllArtists] = useState<Artist[]>([])
  const [loading, setLoading] = useState(true)
  const [inputVal, setInputVal] = useState('')
  const [open, setOpen] = useState(false)
  const [activeIdx, setActiveIdx] = useState(-1)
  const listRef = useRef<HTMLUListElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    getArtists()
      .then(setAllArtists)
      .finally(() => setLoading(false))
  }, [])

  // Sync input display when value changes externally (e.g. cleared)
  useEffect(() => {
    if (value === null) {
      setInputVal('')
    } else {
      const artist = allArtists.find((a) => a.artist_id === value)
      setInputVal(artist?.artist_name ?? value)
    }
  }, [value, allArtists]) // intentionally excludes inputVal

  const filtered = inputVal
    ? allArtists.filter((a) => a.artist_name.toLowerCase().includes(inputVal.toLowerCase()))
    : allArtists

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputVal(e.target.value)
    setOpen(true)
    setActiveIdx(-1)
    if (!e.target.value) onChange(null)
  }

  const selectItem = useCallback(
    (artist: Artist) => {
      setInputVal(artist.artist_name)
      setOpen(false)
      setActiveIdx(-1)
      onChange(artist.artist_id)
    },
    [onChange],
  )

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!open) {
      if (e.key === 'ArrowDown') { setOpen(true); return }
      return
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setActiveIdx((i) => Math.min(i + 1, filtered.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setActiveIdx((i) => Math.max(i - 1, -1))
    } else if (e.key === 'Enter' && activeIdx >= 0) {
      e.preventDefault()
      selectItem(filtered[activeIdx])
    } else if (e.key === 'Escape') {
      setOpen(false)
    }
  }

  // Scroll active item into view
  useEffect(() => {
    if (activeIdx >= 0 && listRef.current) {
      const item = listRef.current.children[activeIdx] as HTMLElement
      item?.scrollIntoView({ block: 'nearest' })
    }
  }, [activeIdx])

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  return (
    <div ref={containerRef} style={{ position: 'relative' }}>
      <input
        type="text"
        disabled={loading}
        placeholder={loading ? 'Loading artists…' : 'Type to search artist…'}
        value={inputVal}
        onChange={handleInputChange}
        onFocus={() => setOpen(true)}
        onKeyDown={handleKeyDown}
        style={inputStyle}
        autoComplete="off"
        spellCheck={false}
      />
      {open && !loading && (
        <ul ref={listRef} style={dropdownStyle}>
          {filtered.length === 0 ? (
            <li style={noMatchStyle}>No artists match</li>
          ) : (
            <>
              {filtered.slice(0, 300).map((artist, i) => (
                <li
                  key={artist.artist_id}
                  style={{
                    ...itemStyle,
                    background: i === activeIdx ? 'var(--accent)' : 'transparent',
                    color: i === activeIdx ? '#fff' : 'var(--text-on-surface)',
                  }}
                  onMouseDown={() => selectItem(artist)}
                  onMouseEnter={() => setActiveIdx(i)}
                >
                  {artist.artist_name}
                </li>
              ))}
              {filtered.length > 300 && (
                <li style={noMatchStyle}>
                  Showing 300 of {filtered.length.toLocaleString()} — keep typing to narrow
                </li>
              )}
            </>
          )}
        </ul>
      )}
    </div>
  )
}

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '8px 10px',
  borderRadius: 6,
  border: '1px solid rgba(225, 224, 221, 0.3)',
  background: 'rgba(241, 242, 240, 0.08)',
  color: 'var(--surface)',
  fontSize: 13,
  fontFamily: 'ui-monospace, Consolas, monospace',
  outline: 'none',
}

const dropdownStyle: React.CSSProperties = {
  position: 'absolute',
  top: 'calc(100% + 4px)',
  left: 0,
  right: 0,
  maxHeight: 260,
  overflowY: 'auto',
  background: 'var(--surface)',
  border: '1px solid var(--muted)',
  borderRadius: 6,
  margin: 0,
  padding: '4px 0',
  listStyle: 'none',
  zIndex: 100,
  boxShadow: '0 8px 24px rgba(0,0,0,0.25)',
}

const itemStyle: React.CSSProperties = {
  padding: '7px 12px',
  fontSize: 12,
  fontFamily: 'ui-monospace, Consolas, monospace',
  cursor: 'pointer',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
}

const noMatchStyle: React.CSSProperties = {
  padding: '10px 12px',
  fontSize: 12,
  color: 'var(--text-muted)',
  fontStyle: 'italic',
}
