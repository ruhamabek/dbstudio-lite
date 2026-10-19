import { ArrowUpRight, Database, Github } from "lucide-react"
import { REPOSITORY_URL } from "@/lib/releases"

export function Navbar() {
  return (
    <header className="site-header">
      <div className="container nav-inner">
        <a href="#top" className="wordmark" aria-label="DBStudio Lite home">
          <Database aria-hidden="true" className="brand-mark" strokeWidth={1.8} />
          <span>DBStudio <span className="brand-lite">Lite</span></span>
        </a>
        <nav aria-label="Main navigation" className="main-nav">
          <a href="#preview">The workspace</a>
          <a href="#features">Features</a>
          <a href={REPOSITORY_URL} target="_blank" rel="noreferrer" className="github-link"><Github size={16} aria-hidden="true" /> GitHub <ArrowUpRight size={13} aria-hidden="true" /></a>
        </nav>
        <a href="#download" className="button button-small button-ink">Download <ArrowUpRight size={15} aria-hidden="true" /></a>
      </div>
    </header>
  )
}
