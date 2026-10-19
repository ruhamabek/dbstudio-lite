import { ArrowDownToLine, ArrowUpRight, Laptop, Monitor, Terminal } from "lucide-react"
import { PLATFORMS, RELEASE_URL } from "@/lib/releases"

const ICONS = { macos: Laptop, windows: Monitor, linux: Terminal }

export function DownloadSection() {
  return (
    <section id="download" className="download-section" aria-labelledby="download-title">
      <div className="container">
        <div className="download-heading"><div><p className="release-label">DBStudio Lite / v0.1.0</p><h2 id="download-title">Make room for<br />a little clarity.</h2></div><p>Your databases are already complex.<br />Working with them doesn’t have to be.<br /><span>Free & open source. Ready for your desktop.</span></p></div>
        <div className="platform-list">{PLATFORMS.map((platform) => {
          const Icon = ICONS[platform.id]
          return <a className="platform-download" href={platform.url} key={platform.id} aria-label={`Download DBStudio Lite for ${platform.name}, ${platform.arch}, ZIP`}><Icon size={26} strokeWidth={1.5} aria-hidden="true" /><div><h3>{platform.name}</h3><p>{platform.arch}</p></div><span className="platform-format">{platform.format}</span><ArrowDownToLine className="platform-arrow" size={21} aria-hidden="true" /></a>
        })}</div>
        <div className="download-footnote"><p>Looking for release notes or another asset?</p><a href={RELEASE_URL} target="_blank" rel="noreferrer">View the release on GitHub <ArrowUpRight size={15} aria-hidden="true" /></a></div>
      </div>
    </section>
  )
}
