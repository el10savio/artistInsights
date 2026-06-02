import { useState } from 'react'

interface Props {
  text: string
}

export default function HelpTooltip({ text }: Props) {
  const [visible, setVisible] = useState(false)

  return (
    <span className="help-wrap">
      <span
        className="help-icon"
        onMouseEnter={() => setVisible(true)}
        onMouseLeave={() => setVisible(false)}
      >
        ?
      </span>
      {visible && <span className="help-bubble">{text}</span>}
    </span>
  )
}
