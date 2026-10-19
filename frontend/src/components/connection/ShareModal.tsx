"use client"

import React, { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { useDatabase } from "@/context/DatabaseContext"
import {
  Share2,
  Copy,
  Check,
  Shield,
  Eye,
  FileText,
  Edit3,
  AlertCircle,
  Unplug,
  Radio,
} from "lucide-react"

interface ShareModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ShareModal({ open, onOpenChange }: ShareModalProps) {
  const { isSharing, activeShareTicket, startSharing, stopSharing } = useDatabase()
  const [permission, setPermission] = useState<"view" | "read" | "update">("read")
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleStartShare = async () => {
    setLoading(true)
    setError(null)
    try {
      await startSharing(permission)
    } catch (err: any) {
      setError(err?.message || String(err))
    } finally {
      setLoading(false)
    }
  }

  const handleStopShare = async () => {
    setLoading(true)
    setError(null)
    try {
      await stopSharing()
    } catch (err: any) {
      setError(err?.message || String(err))
    } finally {
      setLoading(false)
    }
  }

  const handleCopy = async () => {
    if (!activeShareTicket) return
    try {
      await navigator.clipboard.writeText(activeShareTicket)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // fallback
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <div className="flex items-center gap-2">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Share2 className="h-5 w-5" />
            </div>
            <div>
              <DialogTitle className="text-lg">Share Database Session</DialogTitle>
              <DialogDescription className="text-xs text-muted-foreground">
                P2P direct relay to share this database with a teammate without a server.
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        {error && (
          <div className="flex items-center gap-2 rounded-lg border border-destructive/30 bg-destructive/10 p-3 text-xs text-destructive">
            <AlertCircle className="h-4 w-4 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <div className="space-y-4 py-2">
           <div className="flex items-center justify-between rounded-lg border border-border bg-card p-3">
            <div className="flex items-center gap-2.5">
              <Radio className={`h-4 w-4 ${isSharing ? "text-emerald-500 animate-pulse" : "text-muted-foreground"}`} />
              <span className="text-sm font-medium">Sharing Status</span>
            </div>
            {isSharing ? (
              <Badge variant="default" className="gap-1.5 px-2 py-0.5 text-xs bg-emerald-500/10 text-emerald-500 border-emerald-500/20">
                <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                Active Relay
              </Badge>
            ) : (
              <Badge variant="outline" className="text-xs text-muted-foreground">
                Inactive
              </Badge>
            )}
          </div>

          {!isSharing ? (
             <div className="space-y-2.5">
              <label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Granted Permissions
              </label>
              <div className="grid grid-cols-1 gap-2">
                 <button
                  type="button"
                  onClick={() => setPermission("view")}
                  className={`flex items-start gap-3 rounded-lg border p-3 text-left transition-colors ${
                    permission === "view"
                      ? "border-primary bg-primary/5 text-foreground"
                      : "border-border bg-card/50 text-muted-foreground hover:bg-accent/50"
                  }`}
                >
                  <Eye className="h-5 w-5 mt-0.5 shrink-0 text-primary" />
                  <div>
                    <div className="text-sm font-medium text-foreground">View Only</div>
                    <div className="text-xs text-muted-foreground">
                      Schema metadata only (tables & columns). Any raw SQL query is blocked.
                    </div>
                  </div>
                </button>

                {/* Read Only */}
                <button
                  type="button"
                  onClick={() => setPermission("read")}
                  className={`flex items-start gap-3 rounded-lg border p-3 text-left transition-colors ${
                    permission === "read"
                      ? "border-primary bg-primary/5 text-foreground"
                      : "border-border bg-card/50 text-muted-foreground hover:bg-accent/50"
                  }`}
                >
                  <FileText className="h-5 w-5 mt-0.5 shrink-0 text-primary" />
                  <div>
                    <div className="text-sm font-medium text-foreground">Read Only (Recommended)</div>
                    <div className="text-xs text-muted-foreground">
                      Inspect schema + execute SELECT queries. INSERT, UPDATE, DELETE & DDL blocked.
                    </div>
                  </div>
                </button>

                 <button
                  type="button"
                  onClick={() => setPermission("update")}
                  className={`flex items-start gap-3 rounded-lg border p-3 text-left transition-colors ${
                    permission === "update"
                      ? "border-primary bg-primary/5 text-foreground"
                      : "border-border bg-card/50 text-muted-foreground hover:bg-accent/50"
                  }`}
                >
                  <Edit3 className="h-5 w-5 mt-0.5 shrink-0 text-primary" />
                  <div>
                    <div className="text-sm font-medium text-foreground">Read & Write</div>
                    <div className="text-xs text-muted-foreground">
                      Execute SELECT, INSERT, UPDATE, DELETE. Destructive DDL (DROP, TRUNCATE) blocked.
                    </div>
                  </div>
                </button>
              </div>
            </div>
          ) : (
             <div className="space-y-3">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Teammate Share Ticket
                </label>
                <div className="flex gap-2">
                  <Input
                    readOnly
                    value={activeShareTicket || ""}
                    className="font-mono text-xs select-all bg-muted/40"
                  />
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleCopy}
                    className="shrink-0 gap-1.5 px-3"
                  >
                    {copied ? (
                      <>
                        <Check className="h-4 w-4 text-emerald-500" />
                        <span className="text-emerald-500">Copied!</span>
                      </>
                    ) : (
                      <>
                        <Copy className="h-4 w-4" />
                        <span>Copy</span>
                      </>
                    )}
                  </Button>
                </div>
              </div>

              <div className="rounded-lg bg-muted/30 p-3 text-xs text-muted-foreground space-y-1">
                <p className="font-medium text-foreground flex items-center gap-1.5">
                  <Shield className="h-3.5 w-3.5 text-primary" />
                  Zero-Trust P2P Tunnel
                </p>
                <p>
                  Send this ticket to your teammate. They can open DBStudio Lite, choose &quot;Peer Share&quot;,
                  and paste this ticket to connect directly.
                </p>
              </div>
            </div>
          )}
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          {!isSharing ? (
            <Button
              onClick={handleStartShare}
              disabled={loading}
              className="gap-2 w-full sm:w-auto"
            >
              <Share2 className="h-4 w-4" />
              {loading ? "Starting Relay..." : "Create Share Ticket"}
            </Button>
          ) : (
            <Button
              variant="outline"
              onClick={handleStopShare}
              disabled={loading}
              className="gap-2 text-destructive hover:bg-destructive/10 w-full sm:w-auto"
            >
              <Unplug className="h-4 w-4" />
              {loading ? "Stopping..." : "Stop Sharing"}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
