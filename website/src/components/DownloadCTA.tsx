"use client"

import React, { useState, useEffect } from "react"
import { Download, ChevronDown, Check } from "lucide-react"

type Platform = "linux" | "macos" | "windows"

interface PlatformInfo {
  name: string
  label: string
  arch: string
  filename: string
  url: string
}

const PLATFORMS: Record<Platform, PlatformInfo> = {
  linux: {
    name: "Linux",
    label: "Download for Linux",
    arch: "x86_64 (zip)",
    filename: "DBStudio-Lite-Linux-x86_64",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-Linux-x86_64.zip",
  },
  macos: {
    name: "macOS",
    label: "Download for macOS",
    arch: "Universal (Apple Silicon & Intel)",
    filename: "DBStudio-Lite-macOS-Universal.zip",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-macOS-Universal.zip",
  },
  windows: {
    name: "Windows",
    label: "Download for Windows",
    arch: "x64 (.zip)",
    filename: "DBStudio-Lite-Windows-x64.exe",
    url: "https://github.com/ruhamabek/dbstudio-lite/releases/download/v0.1.0/DBStudio-Lite-Windows-x64.zip",
  },
}

export function DownloadCTA() {
  const [selectedPlatform, setSelectedPlatform] = useState<Platform>("linux")
  const [menuOpen, setMenuOpen] = useState(false)

  useEffect(() => {
    if (typeof window === "undefined") return
    const ua = window.navigator.userAgent.toLowerCase()
    if (ua.includes("mac")) {
      setSelectedPlatform("macos")
    } else if (ua.includes("win")) {
      setSelectedPlatform("windows")
    } else {
      setSelectedPlatform("linux")
    }
  }, [])

  const current = PLATFORMS[selectedPlatform]

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="relative inline-flex items-stretch rounded-xl bg-white p-1 shadow-md hover:bg-zinc-100 transition-all">
        <a
          href={current.url}
          className="flex items-center gap-3 px-6 py-3 text-sm font-semibold text-black"
        >
          <Download className="h-4 w-4 shrink-0 text-black" />
          <div className="flex flex-col text-left">
            <span>{current.label}</span>
            <span className="text-[11px] font-normal text-zinc-600 font-mono">
              v0.1.0 • {current.arch}
            </span>
          </div>
        </a>

        <div className="w-[1px] bg-zinc-200 my-1.5" />

         <button
          type="button"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label="Select Operating System"
          className="flex items-center px-3 text-zinc-600 hover:text-black transition-colors rounded-r-lg"
        >
          <ChevronDown className={`h-4 w-4 transition-transform ${menuOpen ? "rotate-180" : ""}`} />
        </button>

         {menuOpen && (
          <div
            className="absolute left-0 right-0 top-full mt-2 rounded-xl border border-white/15 bg-[#121215] p-2 shadow-2xl z-20 backdrop-blur-xl"
            onMouseLeave={() => setMenuOpen(false)}
          >
            <div className="px-2.5 py-1 text-[11px] font-semibold text-zinc-400 uppercase tracking-wider">
              Available Platforms
            </div>
            {(Object.keys(PLATFORMS) as Platform[]).map((key) => {
              const p = PLATFORMS[key]
              const isSelected = key === selectedPlatform
              return (
                <a
                  key={key}
                  href={p.url}
                  onClick={() => {
                    setSelectedPlatform(key)
                    setMenuOpen(false)
                  }}
                  className={`flex items-center justify-between rounded-lg px-3 py-2 text-xs transition-colors ${
                    isSelected
                      ? "bg-white/10 text-white font-semibold"
                      : "text-zinc-300 hover:bg-white/5 hover:text-white"
                  }`}
                >
                  <div className="flex flex-col">
                    <span>{p.name}</span>
                    <span className="text-[10px] text-zinc-400 font-mono">{p.arch}</span>
                  </div>
                  {isSelected && <Check className="h-3.5 w-3.5 text-white" />}
                </a>
              )
            })}
            <div className="mt-1 pt-1 border-t border-white/10">
              <a
                href="https://github.com/ruhamabek/dbstudio-lite/releases/tag/v0.1.0"
                target="_blank"
                rel="noreferrer"
                className="block text-center text-[11px] text-zinc-400 hover:text-white py-1 transition-colors"
              >
                View all release assets on GitHub →
              </a>
            </div>
          </div>
        )}
      </div>

      <p className="text-xs text-zinc-400">
        Free & Open Source under AGPL-3.0 • Standalone Executable (Zero Dependencies)
      </p>
    </div>
  )
}
