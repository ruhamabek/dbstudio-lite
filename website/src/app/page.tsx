import React from "react"
import { Navbar } from "@/components/Navbar"
import { Hero } from "@/components/Hero"
import { PreviewShowcase } from "@/components/PreviewShowcase"
import { Features } from "@/components/Features"
import { DownloadSection } from "@/components/DownloadSection"
import { Footer } from "@/components/Footer"

export default function HomePage() {
   const jsonLd = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "DBStudio Lite",
    operatingSystem: "Linux, macOS, Windows",
    applicationCategory: "DeveloperApplication",
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
    },
    description:
      "A lightweight desktop workspace for PostgreSQL, MySQL, and SQLite. Explore data, inspect schemas, and write SQL.",
    downloadUrl: "https://github.com/ruhamabek/dbstudio-lite/releases/tag/v0.1.0",
    author: {
      "@type": "Person",
      name: "Ruhama",
      url: "https://github.com/ruhamabek",
    },
  }

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />

      <div id="top">
        <a className="skip-link" href="#main-content">Skip to content</a>
        <Navbar />
        <main id="main-content">
         <Hero />
         <PreviewShowcase />
         <Features />
        <DownloadSection />
        </main>
        <Footer />
      </div>
    </>
  )
}
