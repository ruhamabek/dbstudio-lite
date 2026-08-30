"use client"

import React from "react"
import { DatabaseProvider } from "@/context/DatabaseContext"

export function Providers({ children }: { children: React.ReactNode }) {
  return <DatabaseProvider>{children}</DatabaseProvider>
}
