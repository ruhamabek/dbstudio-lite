import type { Metadata } from "next"
import { Analytics } from "@vercel/analytics/react"
import "./globals.css"

export const metadata: Metadata = {
  title: "DBStudio Lite — Fast, Lightweight PostgreSQL Desktop Client",
  description:
    "A modular, lightweight cross-platform PostgreSQL desktop client built with Go, Next.js, and Wails v2. Featuring in-memory speed, Monaco SQL console, and zero external telemetry.",
  keywords: [
    "PostgreSQL",
    "Postgres client",
    "database GUI",
    "Go",
    "Wails",
    "Next.js",
    "DBStudio Lite",
    "SQL editor",
    "Monaco editor",
    "open source",
    "lightweight database client",
  ],
  authors: [{ name: "Ruhama", url: "https://github.com/ruhamabek" }],
  creator: "Ruhama",
  metadataBase: new URL("https://dbstudio-lite.vercel.app"),
  openGraph: {
    title: "DBStudio Lite — Fast, Lightweight PostgreSQL Desktop Client",
    description:
      "A modular, lightweight cross-platform PostgreSQL desktop client built with Go, Next.js, and Wails v2.",
    url: "https://dbstudio-lite.vercel.app",
    siteName: "DBStudio Lite",
    images: [
      {
        url: "/screenshots/2.jpg",
        width: 1200,
        height: 750,
        alt: "DBStudio Lite Table Data Explorer",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "DBStudio Lite — Fast, Lightweight PostgreSQL Desktop Client",
    description:
      "Modular, lightweight cross-platform PostgreSQL desktop client built with Go, Next.js, and Wails v2.",
    images: ["/screenshots/2.jpg"],
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className="dark">
      <body className="min-h-screen bg-[#09090b] text-[#fafafa] antialiased selection:bg-white selection:text-black">
        {children}
        <Analytics />
      </body>
    </html>
  )
}
