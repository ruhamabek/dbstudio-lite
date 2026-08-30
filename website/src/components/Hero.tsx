import React from "react"
import { DownloadCTA } from "./DownloadCTA"
import { Terminal, Shield, Zap } from "lucide-react"

export function Hero() {
  return (
    <section className="relative overflow-hidden pt-20 pb-16 md:pt-28 md:pb-24">
      <div className="pointer-events-none absolute left-1/2 top-0 -translate-x-1/2 h-[450px] w-[800px] bg-radial from-white/5 to-transparent blur-3xl" />
      <div className="relative mx-auto max-w-5xl px-6 text-center">
        <h1 className="text-4xl font-extrabold tracking-tight text-white sm:text-6xl md:text-7xl max-w-4xl mx-auto leading-[1.08]">
          The Minimalist PostgreSQL Client for Developers.
        </h1>
        <p className="mt-6 text-base sm:text-lg text-zinc-400 max-w-2xl mx-auto leading-relaxed">
          Engineered with Go and Wails v2 using strict Test-Driven Development.
          Sub-millisecond execution, embedded Monaco SQL console, and 100% local privacy.
        </p>
        <div className="mt-8">
          <DownloadCTA />
        </div>

        <div className="mt-10 flex flex-wrap items-center justify-center gap-6 text-xs text-zinc-400 font-medium">
          <div className="flex items-center gap-2">
            <Zap className="h-3.5 w-3.5 text-zinc-300" />
            <span>0.004s In-Memory Go Engine</span>
          </div>
          <div className="flex items-center gap-2">
            <Shield className="h-3.5 w-3.5 text-zinc-300" />
            <span>Zero Telemetry or Cloud Tracking</span>
          </div>
          <div className="flex items-center gap-2">
            <Terminal className="h-3.5 w-3.5 text-zinc-300" />
            <span>Monaco Editor Bundled Offline</span>
          </div>
        </div>

        <div className="mt-14 relative rounded-2xl border border-white/10 bg-[#121215] p-2 shadow-2xl backdrop-blur-md">
          <div className="flex items-center justify-between border-b border-white/10 px-4 py-2.5 bg-[#18181b]/50 rounded-t-xl">
            <div className="flex items-center gap-2">
              <div className="h-3 w-3 rounded-full bg-zinc-700" />
              <div className="h-3 w-3 rounded-full bg-zinc-700" />
              <div className="h-3 w-3 rounded-full bg-zinc-700" />
            </div>
            <span className="text-xs font-mono text-zinc-400 font-medium">
              DBStudio Lite — postgres://localhost:5432/production
            </span>
            <div className="w-10" />
          </div>

          <div className="overflow-hidden rounded-b-xl bg-black">
            <img
              src="/screenshots/2.jpg"
              alt="DBStudio Lite Interface Preview"
              className="w-full h-auto object-cover"
              loading="eager"
            />
          </div>
        </div>
      </div>
    </section>
  )
}
