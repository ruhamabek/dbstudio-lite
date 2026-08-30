"use client"

import React, { useState } from "react"
import { useDatabase } from "@/context/DatabaseContext"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { ConnectionModal } from "./ConnectionModal"
import { Database, Plus, Unplug, RefreshCw } from "lucide-react"

export function ConnectionBar() {
  const { connection, disconnect, refreshTables } = useDatabase()
  const [modalOpen, setModalOpen] = useState(false)
  const [isRefreshing, setIsRefreshing] = useState(false)

  const handleRefresh = async () => {
    setIsRefreshing(true)
    await refreshTables()
    setIsRefreshing(false)
  }

  return (
    <>
      <header className="flex h-14 items-center justify-between border-b border-border bg-card px-5 shrink-0">
        {/* Brand & Connection Info */}
        <div className="flex items-center space-x-3.5">
          <div className="flex items-center ">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg text-primary">
              <Database className="h-4.5 w-4.5" />
            </div>
            <span className="font-bold text-base tracking-tight text-foreground">DBStudio Lite</span>
          </div>

          <div className="h-4 w-[1px] bg-border mx-1" />

          {/* Connection Status Badge */}
          {connection.isConnected ? (
            <div className="flex items-center space-x-2.5">
              <Badge variant="default" className="gap-1.5 px-2.5 py-1 text-xs font-semibold">
                <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
                Connected
              </Badge>
              <span className="text-xs text-muted-foreground font-mono truncate max-w-sm">
                {connection.activeDsn}
              </span>
            </div>
          ) : (
            <Badge variant="muted" className="gap-1.5 px-2.5 py-1 text-xs font-medium">
              <span className="h-2 w-2 rounded-full bg-muted-foreground" />
              Disconnected
            </Badge>
          )}
        </div>

        {/* Action Buttons */}
        <div className="flex items-center space-x-2.5">
          {connection.isConnected ? (
            <>
              <Button
                variant="outline"
                size="sm"
                onClick={handleRefresh}
                disabled={isRefreshing}
                className="h-9 px-3.5 gap-1.5 text-sm"
              >
                <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
                Refresh
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={disconnect}
                className="h-9 px-3.5 gap-1.5 text-sm text-destructive hover:bg-destructive/10"
              >
                <Unplug className="h-4 w-4" />
                Disconnect
              </Button>
            </>
          ) : (
            <Button
              size="sm"
              onClick={() => setModalOpen(true)}
              className="h-9 px-4 gap-1.5 text-sm font-semibold shadow-xs"
            >
              <Plus className="h-4 w-4" />
              New Connection
            </Button>
          )}
        </div>
      </header>

      <ConnectionModal open={modalOpen} onOpenChange={setModalOpen} />
    </>
  )
}
