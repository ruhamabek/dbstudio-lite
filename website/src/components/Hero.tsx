import { ArrowDown, ArrowUpRight } from "lucide-react"
import { DownloadCTA } from "./DownloadCTA"
import { REPOSITORY_URL } from "@/lib/releases"

export function Hero() {
  return (
    <section className="hero container" aria-labelledby="hero-title">
      <div className="hero-intro"><span className="status-dot" /> A little lighter. A lot clearer.</div>
      <div className="hero-layout">
        <h1 id="hero-title">Your databases.<br /><span>A clearer view.</span></h1>
        <div className="hero-copy">
          <p>A lightweight desktop workspace for PostgreSQL, MySQL, and SQLite. Explore your data, write SQL, and get back to building.</p>
          <DownloadCTA />
          <div className="hero-note">Free & open source <span aria-hidden="true">·</span> macOS, Windows & Linux</div>
        </div>
      </div>
      <div className="hero-bottom">
        <a href="#preview" className="text-link">Meet your new workspace <ArrowDown size={16} aria-hidden="true" /></a>
        <a href={REPOSITORY_URL} target="_blank" rel="noreferrer" className="text-link hero-source">Built in the open <ArrowUpRight size={15} aria-hidden="true" /></a>
      </div>
    </section>
  )
}
