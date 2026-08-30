import React from "react"
import { Database, Github } from "lucide-react"

export function Footer() {
  return (
    <footer className="border-t border-white/10 bg-[#09090b] py-12 text-zinc-400">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-6 px-6 sm:flex-row">
         <div className="flex items-center gap-3">
          <div className="flex h-7 w-7 items-center justify-center rounded-md bg-white/10 text-white border border-white/10">
            <Database className="h-3.5 w-3.5" />
          </div>
          <span className="text-sm font-bold text-white tracking-tight">
            DBStudio Lite
          </span>
          <span className="text-xs text-zinc-500">
            Released under GNU AGPL-3.0
          </span>
        </div>

         <div className="flex items-center gap-6 text-xs font-medium">
          <a
            href="https://github.com/ruhamabek/dbstudio-lite"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-1.5 hover:text-white transition-colors"
          >
            <Github className="h-3.5 w-3.5" />
            <span>GitHub Repository</span>
          </a>
          <a
            href="https://github.com/ruhamabek/dbstudio-lite/releases"
            target="_blank"
            rel="noreferrer"
            className="hover:text-white transition-colors"
          >
            Releases
          </a>
          <a
            href="https://github.com/ruhamabek/dbstudio-lite/blob/main/LICENSE"
            target="_blank"
            rel="noreferrer"
            className="hover:text-white transition-colors"
          >
            License
          </a>
        </div>
      </div>
    </footer>
  )
}
