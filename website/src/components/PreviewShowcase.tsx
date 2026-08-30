"use client"

import React, { useState } from "react"
import { Table, Terminal, Link, Bookmark } from "lucide-react"

const TABS = [
  {
    id: "data",
    label: "Data Explorer",
    icon: Table,
    title: "Browse & Format Table Records",
    description:
      "Explore rows with intelligent formatting for NULL values, JSON blobs, booleans, and primary key badges. Real-time telemetry shows row counts and query latency.",
    image: "/screenshots/2.jpg",
  },
  {
    id: "console",
    label: "SQL Console",
    icon: Terminal,
    title: "Full-Featured Monaco SQL Editor",
    description:
      "Syntax highlighting, autocompletion, draggable resizable panes, and keyboard shortcuts (Cmd+Enter). Bundled 100% offline with zero CDN dependencies.",
    image: "/screenshots/3.jpg",
  },
  {
    id: "url",
    label: "URL & Parameters",
    icon: Link,
    title: "Flexible Connection Modes",
    description:
      "Connect instantly with a standard PostgreSQL URI or enter discrete parameters. One-click checkbox to store profiles locally.",
    image: "/screenshots/4.png",
  },
  {
    id: "saved",
    label: "Saved Profiles",
    icon: Bookmark,
    title: "1-Click Saved Profiles",
    description:
      "Stored securely with 0600 file permissions on your local disk. Masked passwords in previews and instant one-click connection.",
    image: "/screenshots/1.jpg",
  },
]

export function PreviewShowcase() {
  const [activeTab, setActiveTab] = useState(TABS[0].id)
  const current = TABS.find((t) => t.id === activeTab) || TABS[0]

  return (
    <section id="preview" className="py-20 border-t border-white/10 bg-[#09090b]">
      <div className="mx-auto max-w-6xl px-6">
         <div className="text-center max-w-2xl mx-auto mb-12">
          <h2 className="text-xs font-mono font-bold uppercase tracking-wider text-zinc-400">
            Interface Preview
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Built for Developer Productivity.
          </p>
          <p className="mt-3 text-sm text-zinc-400">
            A clean desktop interface built without visual noise, designed for fast inspection and execution.
          </p>
        </div>

         <div className="flex flex-wrap items-center justify-center gap-2 mb-8">
          {TABS.map((tab) => {
            const Icon = tab.icon
            const isActive = tab.id === activeTab
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 rounded-xl px-4 py-2.5 text-xs font-semibold transition-all ${
                  isActive
                    ? "bg-white text-black shadow-md"
                    : "border border-white/10 bg-white/5 text-zinc-400 hover:bg-white/10 hover:text-white"
                }`}
              >
                <Icon className="h-4 w-4" />
                <span>{tab.label}</span>
              </button>
            )
          })}
        </div>

         <div className="rounded-2xl border border-white/10 bg-[#121215] p-3 md:p-6 shadow-2xl">
          <div className="mb-4 flex flex-col md:flex-row md:items-center justify-between gap-2 px-2">
            <div>
              <h3 className="text-lg font-bold text-white">{current.title}</h3>
              <p className="text-xs text-zinc-400">{current.description}</p>
            </div>
            <span className="text-[11px] font-mono text-zinc-500 shrink-0">
              Active View: {current.label}
            </span>
          </div>

          <div className="overflow-hidden rounded-xl border border-white/10 bg-black">
            <img
              src={current.image}
              alt={current.title}
              className="w-full h-auto object-cover"
              loading="lazy"
            />
          </div>
        </div>
      </div>
    </section>
  )
}
