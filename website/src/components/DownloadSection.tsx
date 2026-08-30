import React from "react"
import { Download, ExternalLink} from "lucide-react"

const BINARIES = [
  {
    os: "Linux",
    arch: "x86_64",
    filename: "DBStudio-Lite-Linux-x86_64",
    format: "Standalone Executable",
    size: "~25 MB",
    compat: "Debian, Ubuntu, Fedora, Arch, Mint",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-Linux-x86_64",
    command: "chmod +x DBStudio-Lite-Linux-x86_64 && ./DBStudio-Lite-Linux-x86_64",
  },
  {
    os: "macOS",
    arch: "Universal (M1/M2/M3/M4 + Intel)",
    filename: "DBStudio-Lite-macOS-Universal.zip",
    format: ".app Bundle (Zipped)",
    size: "~30 MB",
    compat: "macOS 11+ (Big Sur, Monterey, Ventura, Sonoma, Sequoia)",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-macOS-Universal.zip",
    command: "unzip DBStudio-Lite-macOS-Universal.zip",
  },
  {
    os: "Windows",
    arch: "x64",
    filename: "DBStudio-Lite-Windows-x64.exe",
    format: "Portable Executable",
    size: "~24 MB",
    compat: "Windows 10, Windows 11 (64-bit)",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-Windows-x64.exe",
    command: "DBStudio-Lite-Windows-x64.exe",
  },
]

export function DownloadSection() {
  return (
    <section id="download" className="py-24 border-t border-white/10 bg-[#09090b]">
      <div className="mx-auto max-w-6xl px-6">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <h2 className="text-xs font-mono font-bold uppercase tracking-wider text-zinc-400">
            Get DBStudio Lite
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Download for Your Platform
          </p>
          <p className="mt-3 text-sm text-zinc-400">
            Pre-compiled standalone binaries built directly from GitHub Actions release v0.1.0.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {BINARIES.map((b) => (
            <div
              key={b.os}
              className="flex flex-col justify-between rounded-2xl border border-white/10 bg-[#121215] p-6 hover:border-white/20 transition-all"
            >
              <div>
                <div className="flex items-center justify-between mb-4">
                  <span className="text-lg font-bold text-white">{b.os}</span>
                  <span className="rounded-md border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] font-mono text-zinc-400">
                    {b.arch}
                  </span>
                </div>

                <div className="space-y-2 mb-6 text-xs text-zinc-400">
                  <div className="flex justify-between border-b border-white/5 pb-1.5">
                    <span>Format</span>
                    <span className="text-zinc-300 font-mono">{b.format}</span>
                  </div>
                  <div className="flex justify-between border-b border-white/5 pb-1.5">
                    <span>Approx. Size</span>
                    <span className="text-zinc-300 font-mono">{b.size}</span>
                  </div>
                  <div className="flex justify-between pt-0.5">
                    <span>Compatibility</span>
                    <span className="text-zinc-300 text-right max-w-[140px] truncate" title={b.compat}>
                      {b.compat}
                    </span>
                  </div>
                </div>
              </div>

              <div>
                <a
                  href={b.url}
                  className="flex w-full items-center justify-center gap-2 rounded-xl bg-white px-4 py-2.5 text-xs font-bold text-black hover:bg-zinc-200 transition-colors shadow-sm"
                >
                  <Download className="h-4 w-4" />
                  <span>Download {b.os}</span>
                </a>

                {b.command && (
                  <div className="mt-3 rounded-lg bg-black/50 border border-white/5 px-2.5 py-1.5 font-mono text-[10px] text-zinc-400 truncate">
                    $ {b.command}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>

         <div className="mt-12 text-center">
          <a
            href="https://github.com/ruhamabek/dbstudio-lite/releases/tag/v0.1.0"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-2 text-xs font-mono text-zinc-400 hover:text-white transition-colors"
          >
            <span>View checksums and source archives on GitHub Releases</span>
            <ExternalLink className="h-3 w-3" />
          </a>
        </div>
      </div>
    </section>
  )
}
