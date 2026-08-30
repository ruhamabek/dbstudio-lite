import React from "react"
import {
  Cpu,
  ShieldCheck,
  Code2,
  Layers,
  Lock,
  Boxes,
} from "lucide-react"

const FEATURES = [
  {
    icon: Cpu,
    title: "Lightweight Go Core",
    description:
      "Core database inspection and execution engine written in Go with 0.004s test suite latency and minimal RAM consumption.",
  },
  {
    icon: Code2,
    title: "Monaco SQL Console",
    description:
      "Visual Studio Code's editor engine bundled directly into the binary. Syntax highlighting, shortcuts, and draggable split panes.",
  },
  {
    icon: ShieldCheck,
    title: "Zero Telemetry",
    description:
      "Your database queries and schemas never touch the cloud. No analytics, no phone-home servers, and complete offline capability.",
  },
  {
    icon: Layers,
    title: "Schema & PK Inspector",
    description:
      "Inspect tables grouped by schema, discover composite primary keys, detect SQL constraints, and view column data types.",
  },
  {
    icon: Lock,
    title: "0600 Local Profile Security",
    description:
      "Saved connection profiles are stored on local disk under OS configuration directories using user-restricted 0600 file permissions.",
  },
  {
    icon: Boxes,
    title: "Native Desktop Binaries",
    description:
      "Packaged via Wails v2 utilizing the operating system's native webview. No bundled Chromium overhead, producing slim executables.",
  },
]

export function Features() {
  return (
    <section id="features" className="py-24 border-t border-white/10 bg-[#09090b]">
      <div className="mx-auto max-w-6xl px-6">
        {/* Section Title */}
        <div className="text-center max-w-2xl mx-auto mb-16">
          <h2 className="text-xs font-mono font-bold uppercase tracking-wider text-zinc-400">
            Engineered with Precision
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Everything You Need. Nothing You Don't.
          </p>
          <p className="mt-3 text-sm text-zinc-400">
            Designed for database engineers who appreciate lean codebases and instant responsiveness.
          </p>
        </div>

         <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {FEATURES.map((f, i) => {
            const Icon = f.icon
            return (
              <div
                key={i}
                className="group relative rounded-2xl border border-white/10 bg-[#121215] p-6 hover:border-white/20 transition-all hover:shadow-lg"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/5 border border-white/10 text-white group-hover:bg-white/10 transition-colors mb-4">
                  <Icon className="h-5 w-5" />
                </div>
                <h3 className="text-base font-bold text-white mb-2">{f.title}</h3>
                <p className="text-xs leading-relaxed text-zinc-400">{f.description}</p>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
