import { ArrowUpRight, Check, Database, FileCode2, Github, KeyRound } from "lucide-react"
import { REPOSITORY_URL } from "@/lib/releases"

export function Features() {
  return (
    <>
      <div className="database-strip container" aria-label="Supported databases"><p>Three databases.<br /><strong>One familiar workspace.</strong></p><div className="engine"><Database aria-hidden="true" />PostgreSQL</div><div className="engine"><Database aria-hidden="true" />MySQL</div><div className="engine"><FileCode2 aria-hidden="true" />SQLite</div></div>
      <section id="features" className="features-section container" aria-labelledby="features-title">
        <div className="section-heading"><h2 id="features-title">Stay in your flow.<br /><span>Stay close to your data.</span></h2><p>The tools you reach for every day,<br className="desktop-break" /> without a workspace full of distractions.</p></div>
        <div className="feature-layout">
          <article className="query-feature">
            <div className="feature-copy"><h3>A proper editor.<br />Right where you need it.</h3><p>Write SQL with Monaco—the editor behind VS Code. Syntax highlighting, autocomplete, and a resizable results view come built in.</p></div>
            <div className="code-demo" aria-label="Example SQL query">
              <div className="code-title"><span>query.sql</span><span>Example query</span></div>
              <pre><code><span className="sql-comment">-- A little curiosity goes a long way.</span>{"\n"}<span className="sql-keyword">SELECT</span> name, email{"\n"}<span className="sql-keyword">FROM</span> customers{"\n"}<span className="sql-keyword">WHERE</span> status = <span className="sql-string">&apos;active&apos;</span>{"\n"}<span className="sql-keyword">ORDER BY</span> created_at <span className="sql-keyword">DESC</span>;</code></pre>
              <div className="code-footer"><span><Check size={14} aria-hidden="true" /> Bundled offline</span><span><kbd>⌘</kbd> / <kbd>Ctrl</kbd> + <kbd>Enter</kbd> to run</span></div>
            </div>
          </article>
          <div className="feature-details">
            <article><span className="feature-symbol"><Database size={22} aria-hidden="true" /></span><h3>See how it all connects.</h3><p>Explore tables by schema. Inspect column types, defaults, and primary keys without jumping between tools.</p><p className="example-label">Example column</p><div className="schema-example" aria-label="Example column metadata"><span><KeyRound size={13} aria-hidden="true" /> id</span><code>uuid</code><span className="schema-tag">primary key</span></div></article>
            <article><span className="feature-symbol"><KeyRound size={22} aria-hidden="true" /></span><h3>Your connections. Kept local.</h3><p>Save profiles on your own machine with restricted file permissions. Your next session starts with a connection, not a setup form.</p><span className="local-note"><span className="status-dot" /> Stored on your device</span></article>
          </div>
        </div>
      </section>
      <section className="open-source-section" aria-labelledby="open-source-title"><div className="container open-source-inner"><Github size={38} strokeWidth={1.4} aria-hidden="true" /><div><h2 id="open-source-title">Small app. Open book.</h2><p>Built with Go and your system’s native webview. Free to use, with the source right there to explore.</p></div><a className="text-link" href={REPOSITORY_URL} target="_blank" rel="noreferrer">Explore the source <ArrowUpRight size={17} aria-hidden="true" /></a></div></section>
    </>
  )
}
