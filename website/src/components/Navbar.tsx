import React from "react"
import { Database, Github, Download } from "lucide-react"

export function Navbar() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-white/10 bg-[#09090b]/80 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <a href="#" className="flex items-center gap-3 transition-opacity hover:opacity-90">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white/10 text-white border border-white/15 shadow-xs">
            <Database className="h-4 w-4" />
          </div>
          <span className="font-bold tracking-tight text-white text-base">
            DBStudio<span className="text-zinc-400 font-normal ml-1">Lite</span>
          </span>
          <span className="rounded-full border border-white/15 bg-white/5 px-2 py-0.5 text-[11px] font-mono text-zinc-400">
            v0.1.0
          </span>
        </a>

        <nav className="hidden md:flex items-center gap-6 text-sm text-zinc-400 font-medium">
          <a href="#features" className="hover:text-white transition-colors">
            Features
          </a>
          <a href="#preview" className="hover:text-white transition-colors">
            Preview
          </a>
          <a href="#download" className="hover:text-white transition-colors">
            Downloads
          </a>
        </nav>

        <div className="flex items-center gap-3">
          <a
            href="https://github.com/ruhamabek/dbstudio-lite"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-2 rounded-lg border border-white/15 bg-white/5 px-3 py-1.5 text-xs font-medium text-zinc-200 hover:bg-white/10 hover:text-white transition-colors"
          >
            <Github className="h-3.5 w-3.5" />
            <span>GitHub</span>
          </a>
          <a
            href="#download"
            className="flex items-center gap-1.5 rounded-lg bg-white px-3.5 py-1.5 text-xs font-semibold text-black hover:bg-zinc-200 transition-colors shadow-sm"
          >
            <Download className="h-3.5 w-3.5" />
            <span>Get Binary</span>
          </a>
        </div>
      </div>
    </header>
  )
}
