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
import { Input } from "@/components/ui/input"
import { useDatabase } from "@/context/DatabaseContext"
import { db } from "@wailsjs/go/models"
import {
  AlertCircle,
  Database,
  Link,
  Sliders,
  Bookmark,
  Trash2,
} from "lucide-react"

interface ConnectionModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ConnectionModal({ open, onOpenChange }: ConnectionModalProps) {
  const {
    connectWithUrl,
    connectWithConfig,
    connection,
    savedConnections,
    saveCurrentConnection,
    deleteSavedConnection,
  } = useDatabase()

  const [mode, setMode] = useState<"saved" | "url" | "fields">(
    savedConnections.length > 0 ? "saved" : "url"
  )

   const [rawUrl, setRawUrl] = useState(
    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
  )
  const [profileName, setProfileName] = useState("")
  const [shouldSave, setShouldSave] = useState(false)

   const [config, setConfig] = useState<db.Config>(
    new db.Config({
      Host: "localhost",
      Port: 5432,
      User: "postgres",
      Password: "",
      Database: "postgres",
      SSLMode: "disable",
    })
  )

  const handleConnect = async (e: React.FormEvent) => {
    e.preventDefault()
    let success = false
    if (mode === "url") {
      success = await connectWithUrl(rawUrl)
      if (success && shouldSave && profileName.trim()) {
        await saveCurrentConnection(profileName.trim(), rawUrl)
      }
    } else {
      success = await connectWithConfig(config)
    }
    if (success) {
      onOpenChange(false)
    }
  }

  const handleConnectSaved = async (saved: db.SavedConnection) => {
    const success = await connectWithUrl(saved.url)
    if (success) {
      onOpenChange(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent onClose={() => onOpenChange(false)}>
        <form onSubmit={handleConnect} className="space-y-4.5">
          <DialogHeader>
            <div className="flex items-center space-x-2.5">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <Database className="h-4.5 w-4.5" />
              </div>
              <DialogTitle className="text-lg font-bold">Connect to PostgreSQL</DialogTitle>
            </div>
            <DialogDescription className="text-sm">
              Select a saved profile or enter connection credentials.
            </DialogDescription>
          </DialogHeader>

           <div className="flex rounded-lg bg-muted p-1 text-muted-foreground text-sm font-medium">
            <button
              type="button"
              onClick={() => setMode("saved")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "saved"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Bookmark className="h-4 w-4" />
              Saved ({savedConnections.length})
            </button>
            <button
              type="button"
              onClick={() => setMode("url")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "url"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Link className="h-4 w-4" />
              URL
            </button>
            <button
              type="button"
              onClick={() => setMode("fields")}
              className={`flex-1 flex items-center justify-center gap-2 py-2 rounded-md transition-colors ${
                mode === "fields"
                  ? "bg-background text-foreground shadow-xs font-semibold"
                  : "hover:text-foreground"
              }`}
            >
              <Sliders className="h-4 w-4" />
              Parameters
            </button>
          </div>

           {connection.error && (
            <div className="flex items-start space-x-2.5 rounded-lg border border-destructive/20 bg-destructive/10 p-3.5 text-sm text-destructive">
              <AlertCircle className="h-4.5 w-4.5 shrink-0 mt-0.5" />
              <span className="break-all">{connection.error}</span>
            </div>
          )}

           {mode === "saved" && (
            <div className="space-y-2.5 max-h-72 overflow-y-auto pr-1">
              {savedConnections.length > 0 ? (
                savedConnections.map((saved) => (
                  <div
                    key={saved.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-border bg-background hover:border-primary/50 transition-colors"
                  >
                    <div className="flex flex-col min-w-0 pr-3">
                      <span className="text-sm font-semibold text-foreground truncate">
                        {saved.name}
                      </span>
                      <span className="text-xs font-mono text-muted-foreground truncate">
                        {saved.url.replace(/:[^:@]*@/, ":****@")}
                      </span>
                    </div>
                    <div className="flex items-center space-x-2 shrink-0">
                      <Button
                        type="button"
                        size="sm"
                        onClick={() => handleConnectSaved(saved)}
                        disabled={connection.isConnecting}
                        className="h-8 text-xs font-semibold px-3"
                      >
                        Connect
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => deleteSavedConnection(saved.id)}
                        className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="py-10 text-center text-sm text-muted-foreground space-y-1.5">
                  <Bookmark className="h-7 w-7 mx-auto opacity-30" />
                  <p className="font-semibold text-foreground">No saved connections</p>
                  <p className="text-xs max-w-xs mx-auto">
                    Switch to the "URL" tab and check "Save connection" to store profiles here.
                  </p>
                </div>
              )}
            </div>
          )}

           {mode === "url" && (
            <div className="space-y-3.5">
              <div className="space-y-1.5">
                <label className="text-sm font-medium text-foreground">
                  PostgreSQL URI
                </label>
                <Input
                  value={rawUrl}
                  onChange={(e) => setRawUrl(e.target.value)}
                  placeholder="postgres://user:password@localhost:5432/dbname?sslmode=disable"
                  className="font-mono text-xs h-10"
                  required
                />
              </div>

               <div className="rounded-lg border border-border bg-muted/20 p-3.5 space-y-2.5">
                <label className="flex items-center space-x-2.5 cursor-pointer text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={shouldSave}
                    onChange={(e) => setShouldSave(e.target.checked)}
                    className="rounded border-input text-primary focus:ring-primary h-4 w-4"
                  />
                  <span className="text-foreground">
                    Save this connection profile
                  </span>
                </label>
                {shouldSave && (
                  <div className="pt-1">
                    <Input
                      value={profileName}
                      onChange={(e) => setProfileName(e.target.value)}
                      placeholder="e.g. Local Docker, Neon Dev, Staging DB"
                      className="text-sm h-9"
                      required={shouldSave}
                    />
                  </div>
                )}
              </div>
            </div>
          )}

           {mode === "fields" && (
            <div className="grid grid-cols-2 gap-3.5 text-sm">
              <div className="col-span-2 sm:col-span-1 space-y-1.5">
                <label className="font-medium text-muted-foreground">Host</label>
                <Input
                  value={config.Host}
                  onChange={(e) =>
                    setConfig(new db.Config({ ...config, Host: e.target.value }))
                  }
                  placeholder="localhost"
                  className="h-10 text-sm"
                  required
                />
              </div>
              <div className="col-span-2 sm:col-span-1 space-y-1.5">
                <label className="font-medium text-muted-foreground">Port</label>
                <Input
                  type="number"
                  value={config.Port || ""}
                  onChange={(e) =>
                    setConfig(
                      new db.Config({
                        ...config,
                        Port: parseInt(e.target.value, 10) || 5432,
                      })
                    )
                  }
                  placeholder="5432"
                  className="h-10 text-sm"
                  required
                />
              </div>
              <div className="col-span-2 sm:col-span-1 space-y-1.5">
                <label className="font-medium text-muted-foreground">Database</label>
                <Input
                  value={config.Database}
                  onChange={(e) =>
                    setConfig(
                      new db.Config({ ...config, Database: e.target.value })
                    )
                  }
                  placeholder="postgres"
                  className="h-10 text-sm"
                  required
                />
              </div>
              <div className="col-span-2 sm:col-span-1 space-y-1.5">
                <label className="font-medium text-muted-foreground">Username</label>
                <Input
                  value={config.User}
                  onChange={(e) =>
                    setConfig(new db.Config({ ...config, User: e.target.value }))
                  }
                  placeholder="postgres"
                  className="h-10 text-sm"
                  required
                />
              </div>
              <div className="col-span-2 space-y-1.5">
                <label className="font-medium text-muted-foreground">Password</label>
                <Input
                  type="password"
                  value={config.Password}
                  onChange={(e) =>
                    setConfig(
                      new db.Config({ ...config, Password: e.target.value })
                    )
                  }
                  placeholder="••••••••"
                  className="h-10 text-sm"
                />
              </div>
              <div className="col-span-2 space-y-1.5">
                <label className="font-medium text-muted-foreground">SSL Mode</label>
                <select
                  value={config.SSLMode}
                  onChange={(e) =>
                    setConfig(
                      new db.Config({ ...config, SSLMode: e.target.value })
                    )
                  }
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring"
                >
                  <option value="disable">disable</option>
                  <option value="require">require</option>
                  <option value="prefer">prefer</option>
                  <option value="verify-full">verify-full</option>
                </select>
              </div>
            </div>
          )}

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              className="h-10 px-4 text-sm font-medium"
            >
              Cancel
            </Button>
            {mode !== "saved" && (
              <Button
                type="submit"
                disabled={connection.isConnecting}
                className="h-10 px-5 text-sm font-semibold shadow-xs"
              >
                {connection.isConnecting ? "Connecting..." : "Connect"}
              </Button>
            )}
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
