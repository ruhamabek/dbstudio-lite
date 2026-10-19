"use client"

import { useEffect, useState } from "react"
import { Download, ArrowRight } from "lucide-react"
import { PLATFORMS } from "@/lib/releases"

export function DownloadCTA() {
  const [platform, setPlatform] = useState<(typeof PLATFORMS)[number] | null>(null)

  useEffect(() => {
    const ua = navigator.userAgent.toLowerCase()
    if (/android|iphone|ipad|ipod/.test(ua) || (navigator.maxTouchPoints > 1 && /mac/.test(ua))) return
    const id = ua.includes("mac") ? "macos" : ua.includes("win") ? "windows" : ua.includes("linux") ? "linux" : null
    setPlatform(PLATFORMS.find((item) => item.id === id) ?? null)
  }, [])

  return (
    <div className="download-actions">
      <a className="button button-primary" href={platform?.url ?? "#download"}>
        <Download size={18} aria-hidden="true" />
        {platform ? `Download for ${platform.name}` : "Get DBStudio Lite"}
      </a>
      <a className="other-platforms" href="#download">{platform ? "Other platforms" : "See downloads"}<ArrowRight size={14} aria-hidden="true" /></a>
    </div>
  )
}
