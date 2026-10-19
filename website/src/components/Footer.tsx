import { Database } from "lucide-react"
import { REPOSITORY_URL } from "@/lib/releases"

export function Footer() {
  return (
    <footer className="site-footer">
      <div className="footer-main">
        <div className="footer-top">
          <p className="footer-headline">A clearer view starts here.</p>
          <nav className="footer-link-column" aria-label="Product links">
            <h2>Product</h2>
            <a href="#download">Download</a>
            <a href="#preview">Workspace</a>
            <a href="#features">Features</a>
            <a href={`${REPOSITORY_URL}/releases`} target="_blank" rel="noreferrer">Releases</a>
          </nav>
          <nav className="footer-link-column" aria-label="Resources">
            <h2>Resources</h2>
            <a href={`${REPOSITORY_URL}#readme`} target="_blank" rel="noreferrer">Documentation</a>
            <a href={REPOSITORY_URL} target="_blank" rel="noreferrer">Source code</a>
            <a href={`${REPOSITORY_URL}/issues`} target="_blank" rel="noreferrer">Feedback</a>
          </nav>
        </div>
        <a href="#top" className="footer-display" aria-label="DBStudio Lite — back to top">
          <span>DBStudio</span>
          <span className="footer-display-lite">Lite</span>
        </a>
      </div>
      <div className="footer-bottom">
        <a href="#top" className="wordmark"><Database aria-hidden="true" size={24} strokeWidth={1.8} /><span>DBStudio Lite</span></a>
        <p>Free to use. Open by design.</p>
        <a href={`${REPOSITORY_URL}/blob/main/LICENSE`} target="_blank" rel="noreferrer">AGPL-3.0 License</a>
      </div>
    </footer>
  )
}
