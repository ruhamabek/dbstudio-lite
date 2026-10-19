"use client"

import { useRef, useState, type KeyboardEvent } from "react"
import Image from "next/image"
import { Table2, SquareTerminal, Cable, Bookmark, ArrowUpRight } from "lucide-react"

const VIEWS = [
  { id: "data", label: "Explore data", icon: Table2, title: "Less digging. More understanding.", description: "Browse tables, inspect rows, and make sense of your data in one focused view.", image: "/screenshots/2.jpg", alt: "DBStudio Lite data explorer showing database tables in the sidebar and course records in a data grid" },
  { id: "sql", label: "Write SQL", icon: SquareTerminal, title: "From a question to a query.", description: "A familiar Monaco editor, SQL highlighting, and query results right beside your work.", image: "/screenshots/3.jpg", alt: "DBStudio Lite SQL console with the Monaco query editor and query results" },
  { id: "connect", label: "Connect", icon: Cable, title: "Your next connection, made simple.", description: "Use a connection URI or enter your database details, then save a profile for next time.", image: "/screenshots/4.png", alt: "DBStudio Lite connection dialog with connection URL and individual database parameter fields" },
  { id: "profiles", label: "Save profiles", icon: Bookmark, title: "Pick up where you left off.", description: "Keep connection profiles on your machine and reconnect without entering everything again.", image: "/screenshots/1.jpg", alt: "DBStudio Lite saved connection profiles, stored locally for quick reconnection" },
]

export function PreviewShowcase() {
  const [active, setActive] = useState(0)
  const buttons = useRef<(HTMLButtonElement | null)[]>([])

  function onKeyDown(event: KeyboardEvent<HTMLButtonElement>, index: number) {
    let next = index
    if (event.key === "ArrowRight") next = (index + 1) % VIEWS.length
    else if (event.key === "ArrowLeft") next = (index - 1 + VIEWS.length) % VIEWS.length
    else if (event.key === "Home") next = 0
    else if (event.key === "End") next = VIEWS.length - 1
    else return
    event.preventDefault()
    setActive(next)
    buttons.current[next]?.focus()
  }

  return (
    <section id="preview" className="preview-section" aria-label="Explore the DBStudio Lite workspace">
      <div className="container">
        <div className="preview-toolbar">
          <div className="preview-tabs" role="tablist" aria-label="Workspace previews">
            {VIEWS.map((view, index) => {
              const Icon = view.icon
              return <button key={view.id} id={`tab-${view.id}`} ref={(node) => { buttons.current[index] = node }} role="tab" aria-selected={active === index} aria-controls={`panel-${view.id}`} tabIndex={active === index ? 0 : -1} onKeyDown={(event) => onKeyDown(event, index)} onClick={() => setActive(index)}><Icon size={16} aria-hidden="true" />{view.label}</button>
            })}
          </div>
          <span className="preview-version">A look inside · v0.1.0</span>
        </div>
        {VIEWS.map((current, index) => (
        <div key={current.id} id={`panel-${current.id}`} role="tabpanel" aria-labelledby={`tab-${current.id}`} tabIndex={0} hidden={active !== index} className="preview-panel">
          <div className="app-window">
            <div className="window-bar"><div className="window-dots" aria-hidden="true"><i /><i /><i /></div><span>DBStudio Lite</span><span className="window-local">Desktop workspace</span></div>
            <a href={current.image} target="_blank" rel="noreferrer" className={`screenshot-link screenshot-${current.id}`} aria-label={`Open full-size screenshot: ${current.label}`}>
              <Image src={current.image} alt={current.alt} width={1280} height={744} priority={index === 0} sizes="(max-width: 1200px) 92vw, 1120px" className="app-screenshot" />
              <span className="expand-preview">View full size <ArrowUpRight size={14} aria-hidden="true" /></span>
            </a>
          </div>
          <div className="preview-caption"><h2>{current.title}</h2><p>{current.description}</p></div>
        </div>
        ))}
      </div>
    </section>
  )
}
